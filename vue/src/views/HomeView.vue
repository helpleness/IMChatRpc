<template>
  <div class="center-container">
    <!-- 顶部导航栏，包含导航链接和登录按钮 -->
    <div class="top-bar">
      <nav class="navbar">
        <ul class="nav-links">
          <li><a href="/">主页</a></li>
          <li><a href="/about">关于我们</a></li>
          <li><a href="/contact">联系我们</a></li>
        </ul>
      </nav>
      <div class="login-area">
        <div v-if="!chatStore.isLoggedIn" class="login-button">
          <Button label="登录" @click="login" />
        </div>
        <div v-else class="user-avatar">
          <Avatar image="https://www.w3schools.com/w3images/avatar2.png" shape="circle" />
          <span class="username">{{ chatStore.username }}</span>
        </div>
      </div>
    </div>

    <!-- 功能区：只有登录后可用 -->
    <Card v-if="chatStore.isLoggedIn" class="card-layout">
      <template #content>
        <Splitter style="height: 600px">

          <SplitterPanel :size="25" :minSize="25">
            <!-- 侧边栏 -->
            <ScrollPanel style="width: 100%; height: 100%" class="p-3">
              <Button icon="pi pi-comments" label="聊天室" class="wechat-button chat-button" @click="gochatroom" />
              <Button icon="pi pi-user-plus" label="通讯录" class="wechat-button friends-button" @click="goaddress"/>
            </ScrollPanel>
          </SplitterPanel>
          <!-- 左侧用户列表 -->
          <SplitterPanel class="user-panel" :size="45" :minSize="45" :maxSize="45">
            <ScrollPanel style="width: 100%; height: 100%" class="p-3">
              <div class="search-and-controls">
                <div class="search-bar">
                  <i class="pi pi-search search-icon"></i>
                  <InputText v-model="searchQuery" placeholder="搜索" class="search-input" />
                </div>
                <!-- SpeedDial 功能按钮 -->
                <div class="speeddial-bar">
                  <SpeedDial :model="items" direction="down" :style="{ position: 'absolute' }" :buttonProps="{ severity: 'warn', rounded: true }"/>
                </div>
              </div>

              <h3>在线好友</h3>
              <ul class="user-list">
                <li v-for="(user, index) in chatStore.onlineusers" :key="index" class="mb-2 user-item" @click="selectUser(user)">
                  {{ user }}
                </li>
              </ul>
            </ScrollPanel>
          </SplitterPanel>

          <!-- 右侧聊天窗口 -->
          <SplitterPanel class="chat-panel" :size="150">
            <dev v-if="chatStore.selectedUser !== '' ">
              <!-- 当前聊天对象名称显示 -->
              <div v-if="chatStore.selectedUser"  style="width: 100%; height: 10%;">
                <h3>与 {{ chatStore.selectedUser }} 的聊天</h3>
              </div>
              <!-- 消息展示区 -->
              <ScrollPanel style="width: 100%; height: 60%;" class="message-panel" ref="messagePanel" >
                <div v-for="(msg, index) in chatStore.messages" :key="index" class="mb-3 message-item">
                  <strong>{{ msg.user }}:</strong>
                  <span>{{ msg.content }}</span>
                </div>
              </ScrollPanel>

              <ScrollPanel style="width: 100%; height: 30%;" class="message-panel">
                <!-- 消息输入框和控制按钮 -->
                <Toast />
                <!-- 控制按钮 -->
                <div class="control-buttons flex">
                  <Button icon="pi pi-phone" class="icon-button voice-call-button" @click="startVoiceCall" />
                  <Button icon="pi pi-video" class="icon-button video-call-button" @click="startVideoCall" />
                  <Button icon="pi pi-microphone" class="icon-button audio-message-button" @click="recordAudioMessage" />
                  <FileUpload mode="basic" name="demo[]" url="/api/upload" accept="image/*" :maxFileSize="1000000" @upload="onUpload" :auto="true" chooseLabel="" class="icon-button upload-icon" icon="pi pi-image" />
                </div>
                <div class="flex message-input-container align-center">
                  <!-- 消息输入框 -->
                  <InputText v-model="chatStore.newMessage" placeholder="输入消息..." class="w-full message-input" @keyup.enter="sendNewMessage" />
                  <!-- 发送按钮 -->
                  <Button icon="pi pi-send" class="icon-button send-button" @click="sendNewMessage" />
                </div>
              </ScrollPanel>
            </dev>

          </SplitterPanel>
        </Splitter>
      </template>
    </Card>

    <!-- 未登录提示 -->
    <Card v-else class="login-prompt-card">
      <template #content>
        <h2>请先登录以使用此页面的功能。</h2>
      </template>
    </Card>

    <!-- 视频通话窗口 -->
    <div v-if="chatStore.isVideoCalling" class="video-container">
      <div class="video-wrapper">
        <video id="localVideo" autoplay playsinline></video>
        <video id="remoteVideo" autoplay playsinline></video>
      </div>
      <Button label="结束通话" @click="endVideoCall" class="end-call-button" />
    </div>
  </div>
</template>

<script>
import { onMounted, onBeforeUnmount } from 'vue';
import { watch, nextTick } from 'vue';
import { useChatStore } from "@/stores/chatStore";
import router from "@/router";
import {useWebSocketStore} from "@/stores/WebSocket.js";
import {computed, ref} from "vue";
import {use_send_message} from "@/stores/rpc/rpcmessage.js"
import {use_send_group_join_request} from "@/stores/rpc/rpcgroupjoin.js";
import {use_handle_friend_request} from "@/stores/rpc/rpchandle_friend.js";
import {use_create_group} from "@/stores/rpc/create_group.js";
import {use_receive_message_stream} from "@/stores/rpc/receive_message_stream.js";
import {use_send_friend_request} from "@/stores/rpc/rpcfriends.js";
import {use_handle_group_join_request} from "@/stores/rpc/handle_group_join_request.js"
import { proto } from '@/stores/rpc/IMChat_grpc_web_pb.js';
const ChatServiceClient = proto.chat.ChatServiceClient// 导入生成的 ChatServiceClient
import axios from "axios";

export default {
  setup() {
    //
    // - `POST /v1/send_message`：发送消息
    // - `POST /v1/send_friend_request`：发送好友申请
    // - `POST /v1/send_group_join_request`：发送群聊加入申请
    // - `POST /v1/handle_friend_request`：处理好友申请
    // - `POST /v1/handle_group_join_request`：处理群聊加入申请
    // - `POST /v1/create_group`：创建群聊
    // - `GET /v1/receive_message_stream/{user_id}`：接收实时消息流
    const chatStore = useChatStore();
    const websocketStore = useWebSocketStore();
    const send_message =use_send_message();
    const send_friend_request =use_send_friend_request();
    const send_group_join_request=use_send_group_join_request();
    const handle_friend_request=use_handle_friend_request();
    const handle_group_join_request=use_handle_group_join_request();
    const create_group=use_create_group();
    const receive_message_stream=use_receive_message_stream();
    const newMessage = ref('');
    const messages = ref([]); // 存储聊天消息
    const searchQuery = ref('');
    const client = new ChatServiceClient('http://localhost/rpc', null, null);
    chatStore.client_connect(client);

    //const selectedUser = ref(null);
    const items = ref([
      {
        label: '发起群聊', // 修改为发起群聊
        icon: 'pi pi-users', // 使用 pi-users 图标表示群聊
        command: () => {
          createGroup();
        }
      },
      {
        label: '添加朋友', // 修改为添加朋友
        icon: 'pi pi-user-plus', // 使用 pi-user-plus 图标表示添加朋友
        command: () => {
          addFriend();
        }
      }
    ]);
    const messagePanel = ref(null);
    // 过滤用户列表，匹配搜索查询
    const filteredUsers = computed(() => {
      return chatStore.users.filter(user =>
          user.toLowerCase().includes(searchQuery.value.toLowerCase())
      );
    });

    onMounted(() => {
      receive_message_stream.startReceivingMessages(client);
    });

    onBeforeUnmount(() => {
      receive_message_stream.stopReceivingMessages();
    });

    // 选择聊天的用户
    const selectUser = (user) => {
      chatStore.selected_User(user);
      chatStore.loadMessagesForUser(user); // 加载选中用户的聊天记录
    };

    // 加好友功能
    const addFriend = () => {

      alert('加好友功能开发中...');
    };

    // 建群功能
    const createGroup = () => {
      alert('建群功能开发中...');
    };
    const gochatroom =() =>{
      // router.push({name:"home"})
    };
    const goaddress =() =>{
      router.push({ name: "address"})
    }
    // Login action
    const login = () => {
      router.push({ name: "login" });
    };

    const onUpload = () => {
      this.$toast.add({
        severity: 'info',
        summary: 'Success',
        detail: 'File Uploaded',
        life: 3000,
      });
    };

    // Send a new message
    const sendNewMessage = () => {
      if (chatStore.newMessage.trim()) {
        send_message.sendMessage(chatStore.newMessage,chatStore.username,chatStore.selectedUser,chatStore.client)
        // 发送成功，调用 chatStore 处理消息
        console.log("success:",JSON.stringify(send_message.success))

        // websocketStore.sendMessage(chatStore.messages);
      }
    };

    const recordAudioMessage = () => {
      if (!navigator.mediaDevices) {
        alert('此设备不支持音频录制');
        return;
      }
      navigator.mediaDevices
          .getUserMedia({ audio: true })
          .then((stream) => {
            const mediaRecorder = new MediaRecorder(stream);
            const audioChunks = [];

            mediaRecorder.ondataavailable = (event) => {
              audioChunks.push(event.data);
            };

            mediaRecorder.onstop = () => {
              const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
              const audioUrl = URL.createObjectURL(audioBlob);
              messages.value.push({ user: '我', content: `语音消息: ${audioUrl}` });
              const audio = new Audio(audioUrl);
              audio.play(); // 播放录制的音频
            };

            mediaRecorder.start();
            setTimeout(() => mediaRecorder.stop(), 5000); // 录制5秒
          })
          .catch((error) => {
            console.error('录音失败', error);
          });
    };

    // Start voice call (future implementation)
    const startVoiceCall = () => {
      alert("语音通话功能开发中...");
    };


    // Start video call
    const startVideoCall = () => {
      chatStore.isVideoCalling = true;
      navigator.mediaDevices
          .getUserMedia({ video: true, audio: true })
          .then((stream) => {
            const localVideo = document.getElementById("localVideo");
            localVideo.srcObject = stream;
          })
          .catch((error) => {
            console.error("无法访问摄像头/麦克风", error);
          });
    };

    // End video call
    const endVideoCall = () => {
      chatStore.isVideoCalling = false;
      const localVideo = document.getElementById("localVideo");
      const stream = localVideo.srcObject;
      if (stream) {
        const tracks = stream.getTracks();
        tracks.forEach((track) => track.stop());
      }
    };

    watch(
        () => chatStore.messages,
        async () => {
          await nextTick();
          if (messagePanel.value) {
            if (typeof messagePanel.value.scrollToBottom === 'function') {
              messagePanel.value.scrollToBottom();
            } else {
              const scrollContent = messagePanel.value.$el.querySelector('.p-scrollpanel-content');
              if (scrollContent) {
                scrollContent.scrollTop = scrollContent.scrollHeight;
              }
            }
          }
        },
        { deep: true }
    );
    return {
      chatStore,
      newMessage,
      messages,
      searchQuery,
      filteredUsers,
      login,
      sendNewMessage,
      startVoiceCall,
      startVideoCall,
      endVideoCall,
      recordAudioMessage,
      onUpload,
      selectUser,
      addFriend,
      createGroup,
      items,
      gochatroom,
      goaddress,
      messagePanel,

    };
  },
};
</script>

<style scoped>

.wechat-button {
  height: auto;
  background-color: #1aad19; /* 微信绿色 */
  color: white;
  border: none;
  margin: 5px 0;
  border-radius: 0 10px 10px 0; /* 圆角仅在左侧 */
}

.wechat-button:hover {
  background-color: #128c7e; /* 深绿色 */
}

.search-and-controls {
  display: flex;
  justify-content: space-between;
  padding-top: 0.2rem;
  padding-bottom: 0.2rem;
  width: 100%;
  height: 8%;
}

.search-bar {
  display: flex;
  width: 70%;
}

.search-icon {
  font-size: 1.5rem;
  color: #999;
  padding: 0.5rem 0.5rem;
}

.search-input {
  width: 100%;
  padding: 0.5rem 0.5rem;
  border-radius: 20px;
  border: 1px solid #ddd;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  font-size: 1rem;
}

.speeddial-bar {
  display: flex;
  width: 30%;
  padding: 0 1rem 0 1rem;

}

.chat-header {
  background-color: #f1f1f1;
  align-self: center;
}

.top-bar {
  width: 100%;
  padding: 1rem;
  background-color: #ffffff;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: absolute;
  top: 0;
}

.navbar {
  display: flex;
  align-items: center;
}

.nav-links {
  list-style: none;
  display: flex;
  margin: 0;
  padding: 0;
}

.nav-links li {
  margin-right: 2rem;
}

.nav-links a {
  text-decoration: none;
  color: #333;
  font-size: 1.2rem;
  transition: color 0.3s ease;
}

.nav-links a:hover {
  color: #42a5f5; /* 修改导航链接的 hover 颜色 */
}

.login-area {
  display: flex;
  align-items: center;
}

.login-button {
  margin-right: 1rem;
}

.username {
  margin-left: 0.5rem;
  font-size: 1.2rem;
  color: #333;
}

.login-prompt-card {
  text-align: center;
}

.message-panel {
  background-color: #f9f9f9;
  border-radius: 10px;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1rem;
}

.center-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #e0f7fa, #f9f9f9); /* 更柔和的渐变背景 */
}

.top-bar {
  width: 100%;
  padding: 1rem;
  background-color: #ffffff;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: absolute;
  top: 0;
}

.navbar {
  display: flex;
  align-items: center;
}

.nav-links {
  list-style: none;
  display: flex;
  margin: 0;
  padding: 0;
}

.nav-links li {
  margin-right: 2rem;
}

.nav-links a {
  text-decoration: none;
  color: #333;
  font-size: 1.2rem;
  transition: color 0.3s ease;
}

.nav-links a:hover {
  color: #42a5f5; /* 修改导航链接的 hover 颜色 */
}

.login-area {
  display: flex;
  align-items: center;
}

.login-button {
  margin-right: 1rem;
}

.user-avatar {
  display: flex;
  align-items: center;
}

.username {
  margin-left: 0.5rem;
  font-size: 1.2rem;
  color: #333;
}

.card-layout,
.login-prompt-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 2rem;
  border-radius: 15px;
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1); /* 更加柔和的阴影 */
  background: #ffffff;
  margin-top: 5rem;
  width: 70%;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.card-layout:hover {
  transform: translateY(-5px); /* 提升交互感，卡片悬浮效果 */
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.login-prompt-card {
  text-align: center;
}

.user-panel {
  background-color: #f5f5f5;
  border-right: 2px solid #e0e0e0;
}

.chat-panel {
  padding: 1.5rem;
  background: #ffffff;
  border-radius: 15px;

}

.user-list {
  list-style: none;
  padding: 0;
}

.user-item {
  padding: 0.5rem;
  background-color: #e1f5fe;
  border-radius: 10px;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
}

.user-item:hover {
  background-color: #b3e5fc;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.1);
}

.message-item {
  padding: 0.8rem;
  background-color: #f1f8e9;
  border-radius: 10px;
  margin-bottom: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.chat-controls-container {
  display: flex;
  align-items: center;
  padding: 1rem;
  position: absolute;
  bottom: 0;
  width: 50%;
  background-color: #fff;
  box-shadow: 0 -1px 5px rgba(0, 0, 0, 0.1);
}

.message-input-container {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 15px; /* 增加间距 */
}

.message-input {
  flex-grow: 1;
  border-radius: 20px;
  padding: 0.75rem 1rem;
  border: 1px solid #ddd;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05); /* 添加阴影以提高视觉层次 */
  font-size: 16px;
  background: #f9f9f9;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

.message-input:focus {
  background-color: #ffffff;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 焦点时更明显的阴影 */
}

.send-button {
  background-color: #42a5f5;
  color: white;
  border-radius: 50%;
  width: 55px;
  height: 55px;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

.send-button:hover {
  background-color: #1e88e5;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.1); /* 增加 hover 阴影效果 */
}

.control-buttons {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0.3rem;
}

.icon-button {
  width: 45px;
  height: 45px;
  border-radius: 50%;
  background-color: #42a5f5;
  color: white;
  font-size: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

.icon-button:hover {
  background-color: #1e88e5;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.1); /* 按钮 hover 时的阴影 */
}

.upload-icon {
  background-color: #ffeb3b;
}

.upload-icon:hover {
  background-color: #fdd835;
}

.video-container {
  position: fixed;
  top: 10%;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  padding: 1rem;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.video-wrapper {
  display: flex;
  justify-content: space-between;
  width: 100%;
  margin-bottom: 1rem;
}

video {
  width: 300px;
  height: 200px;
  background-color: #000;
  border-radius: 10px;
  margin: 0 10px;
}

.end-call-button {
  background-color: #e53935;
  color: #fff;
  padding: 0.75rem 2rem;
  font-size: 1rem;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.end-call-button:hover {
  background-color: #d32f2f;
}

</style>

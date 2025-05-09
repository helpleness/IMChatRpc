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
              <Button icon="pi pi-user-plus" label="新的朋友" @click="newfirends" class="wechat-button chat-button"/>
              <h3>好友列表</h3>
              <ul class="user-list">
                <li v-for="(user, index) in chatStore.allusers" :key="index" class="mb-2 user-item" @click="selectUser(user)">
                  {{ user }}
                </li>
              </ul>
            </ScrollPanel>
          </SplitterPanel>

          <!-- 右侧聊天窗口 -->
          <SplitterPanel class="chat-panel" :size="150">
            <div v-if="chatStore.selectedUser !== '' ">
              <!-- 当前聊天对象名称显示 -->
              <div class="user-show">
                <Avatar image="https://www.w3schools.com/w3images/avatar2.png" shape="circle" style="width: 100px; height: 100px; border-radius: 50%;" />
                <span class="usernameshow">{{ chatStore.selectedUser }}</span>
              </div>
              <div class="control-buttonsshow flex" >
                <Button icon="pi pi-phone" class="icon-button voice-call-button" @click="startVoiceCall" />
                <Button icon="pi pi-video" class="icon-button video-call-button" @click="startVideoCall" />
                <Button icon="pi pi-comments" class="icon-button voice-call-button" @click="gochatroom" />
              </div>
            </div>
          </SplitterPanel>
        </Splitter>
      </template>
    </Card>

    <!-- 未登录提示 -->
    <Card v-else class="login-prompt-card">
      <h2>请先登录以使用此页面的功能。</h2>
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
import { useChatStore } from "@/stores/chatStore";
import router from "@/router";
import {useWebSocketStore} from "@/stores/WebSocket.js";
import {computed, ref} from "vue";

export default {
  setup() {
    const chatStore = useChatStore();
    const websocketStore = useWebSocketStore();
    const newMessage = ref('');
    const messages = ref([]); // 存储聊天消息
    const searchQuery = ref('');
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
    // 过滤用户列表，匹配搜索查询
    const filteredUsers = computed(() => {
      return chatStore.users.filter(user =>
          user.toLowerCase().includes(searchQuery.value.toLowerCase())
      );
    });

    // 选择聊天的用户
    const selectUser = (user) => {
      chatStore.selected_User(user);
      // chatStore.loadMessagesForUser(user); // 加载选中用户的聊天记录
    };

    // 建群功能
    const createGroup = () => {
      alert('建群功能开发中...');
    };
    // 加好友功能
    const addFriend = () => {
      alert('加好友功能开发中...');
    };

    const newfirends =() =>{
      alert('加好友功能开发中...');
    };

    const gochatroom =() =>{
      router.push({name:"home"})
    };
    const goaddress =() =>{
      //router.push({name:"address"})
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


    return {
      chatStore,
      newMessage,
      messages,
      searchQuery,
      filteredUsers,
      login,
      startVoiceCall,
      startVideoCall,
      endVideoCall,
      onUpload,
      selectUser,
      addFriend,
      createGroup,
      items,
      gochatroom,
      goaddress,
      newfirends,

    };
  },
};
</script>

<style scoped>


.user-show {
  display: flex;
  align-items: center;
  justify-content: center; /* 居中显示 */
  padding: 10rem 12rem 0 12rem;
}

.user-show .pi {
  width: 80px; /* 头像宽度 */
  height: 80px; /* 头像高度 */
  border-radius: 50%; /* 圆形头像 */
}

.user-show .usernameshow {
  font-size: 3rem; /* 用户名字体大小 */
  margin-left: 1rem; /* 头像和用户名之间的间隔 */
}


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

.control-buttonsshow {
  display: flex;
  align-items: center;
  gap: 50px;
  padding: 3rem;
  justify-content: center;
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

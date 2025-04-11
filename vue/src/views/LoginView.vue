<template>
  <div class="center-container" style="text-align: center;">
    <Card class="card-layout" style="width: 30rem; overflow: hidden">
      <template #header>
        <img alt="user header" src="https://primefaces.org/cdn/primevue/images/usercard.png" />
      </template>
      <template #title>登录</template>
      <template #content>
        <!-- 使用 gap 控制各个模块之间的间距 -->
        <div class="form-group">
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-user"></i>
            </InputGroupAddon>
            <InputText v-model="Username_value" placeholder="Username" />
          </InputGroup>
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-key"></i>
            </InputGroupAddon>
            <Password v-model="Password_value" placeholder="Password" toggleMask />
          </InputGroup>
          <div class="flex items-center">
            <Checkbox v-model="checked" inputId="ingredient" name="checked" value="记住我" />
            <label for="ingredient" class="ml-2"> 记住我 </label>
          </div>
        </div>
      </template>
      <template #footer>
        <div class="flex gap-4 mt-1">
          <Button label="登录" @click="login" size="large" class="w-full" />
        </div>
      </template>
    </Card>
  </div>
</template>

<script>
import { ref } from 'vue'
import Password from "primevue/password";
import {useUserStore} from "@/stores/user.js";
import {useChatStore} from "@/stores/chatStore.js";
import {useWebSocketStore} from "@/stores/WebSocket.js";
export default {
  data(){
    return {
      checked: false,
      Password_value: "",
      Username_value: ""
    };
  },
  methods:{
    login(){
      console.log(this.Username_value,this.Password_value)
      const userStore  = useUserStore();
      const chatstore = useChatStore();
      const websocketStore = useWebSocketStore();
      userStore.login(this.Username_value,this.Password_value).then((res) =>{
        //登录成功 推送
        chatstore.login(userStore.info.ID)
        websocketStore.connect()
        chatstore.getusername( userStore.info.username )
        this.$router.push({name:"home"})
      }).catch((err)=>{
        //登录失败 弹窗
        this.$toast.add({ severity: 'error', summary: 'Error Message', detail: err, life: 3000 });
      })
    }
  }
}
</script>

<style scoped>
.center-container {
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  height: 100vh; /* 占满整个视口高度 */
}

/* 使用 flex 布局并设置 gap 来控制上下间距 */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 1.5rem; /* 控制每个元素之间的上下间距 */
}

.ml-2 {
  margin-left: 0.5rem; /* 调整复选框和标签之间的水平间距 */
}

.card-layout {
  display: flex;
  flex-direction: column;
  justify-content: space-between; /* 在卡片内部均匀分布内容 */
}
</style>

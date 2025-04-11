<template>
  <div class="center-container" style="text-align: center;">
    <Toast />
    <Card class="card-layout" style="width: 30rem; overflow: hidden">
      <template #header>
        <img alt="user header" src="https://primefaces.org/cdn/primevue/images/usercard.png" />
      </template>
      <template #title>注册</template>
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
        </div>
      </template>
      <template #footer>
        <div class="flex gap-4 mt-1">
          <Button label="注册" @click="register" size="large" class="w-full" />
        </div>
      </template>
    </Card>
  </div>
</template>

<script>
import { ref } from 'vue'
import Password from "primevue/password";
import {useUserStore} from "@/stores/user.js";
export default {
  data(){
    return {
      Password_value: "",
      Username_value: ""
    };
  },
  methods:{
    register(){
      console.log(this.Username_value,this.Password_value)
      const userStore  = useUserStore();
      userStore.register(this.Username_value,this.Password_value).then((res) =>{
        //注册成功 推送
        this.$router.push({name:"dashboard"})
      }).catch((err)=>{
        //注册失败 弹窗
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

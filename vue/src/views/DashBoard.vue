<template>
  <div>
    {{userinfo}}
  </div>
</template>

<script>
import {useUserStore} from "@/stores/user.js";

export default {
  data(){
    return{
      userinfo:""
    }
  },
  mounted() {
    this.getUserInfo()
  },
  methods:{
    getUserInfo(){
      const userStore = useUserStore();
      this.userinfo = userStore.info
      console.log(this.userinfo )
    }
  }
}
</script>

<style scoped>

</style>

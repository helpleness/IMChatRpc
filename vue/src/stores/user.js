import {defineStore} from "pinia";
import qs from 'qs';
import storage from "@/service/storage.js";
import userService from '@/service/user.js'
export const useUserStore = defineStore("user",{
    state:() => ({
        token: storage.get(storage.USER_TOKEN),
            info: storage.get(storage.USER_INFO)
            ? JSON.parse(storage.get(storage.USER_INFO))
            :null,
    }),
    actions:{
        SET_TOKEN(token){
            storage.set(storage.USER_TOKEN,token)
            this.token =token
        },
        SET_USERINFO(userinfo){
            storage.set(storage.USER_INFO,JSON.stringify(userinfo))
            this.userInfo =userinfo
        },
        register(username,  password){
            return new Promise((resolve,reject) => {
                userService.register( {username , password})
                    .then(
                        (res)=>{
                            this.SET_TOKEN(res.data.data.token)
                            //res.data={
                            // 		"code": 200,
                            // 		"data": gin.H{
                            // 			"token": token,
                            // 		},
                            // 		"msg": "login success",
                            // 	}
                            userService.info().then((res)=>{
                                this.SET_USERINFO(res.data.data)
                                //{
                                // 		"code": "200",
                                // 		"data": gin.H{
                                // 			"ID":       u.ID,
                                // 			"username": u.Username,
                                // 		},
                                // 		"msg": "注册成功",
                                // 	}
                                resolve(res)
                            })
                        })
                    .catch((err) => {
                        reject(err)
                    })
            })
        },
        login(username,  password){
            return new Promise((resolve,reject) => {
                userService.login( {username , password})
                    .then(
                        (res)=>{
                            this.SET_TOKEN(res.data.data.token)
                            //res.data={
                            // 		"code": 200,
                            // 		"data": gin.H{
                            // 			"token": token,
                            // 		},
                            // 		"msg": "login success",
                            // 	}
                            userService.info().then((res)=>{
                                this.SET_USERINFO(res.data.data)
                                //{
                                // 		"code": "200",
                                // 		"data": gin.H{
                                // 			"ID":       u.ID,
                                // 			"username": u.Username,
                                // 		},
                                // 		"msg": "注册成功",
                                // 	}
                                resolve(res)
                            })
                        })
                    .catch((err) => {
                        reject(err)
                    })
            })
        }

    }
})

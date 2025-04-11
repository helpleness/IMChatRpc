import axios from 'axios'
import storage from "@/service/storage.js";

const service =axios.create({
    baseURL: '/api/',
    timeout: 1000 * 5,
})

service.interceptors.request.use((config)=>{
    Object.assign(config.headers,{Authorization: `Bearer ${storage.get(storage.USER_TOKEN)}`})
    return config
},(error)=>(Promise.reject(error)))

export default service;

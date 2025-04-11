import request from '@/util/request.js'
import qs from "qs";

const register =({username,password}) =>(request.post('/register',qs.stringify({username , password})))

const login =({username,password})=> (request.post('/login',qs.stringify({username , password})))

const info =() => request.get('/userinfo')

export default {
    register,
    login,
    info
}

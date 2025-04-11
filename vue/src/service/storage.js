//本地缓存服务

const PREFIX = 'gwq_';
//用户头
const USER_PREFIX = `${PREFIX}user_`;

//要存的东西 ${USER_PREFIX}字符串拼接
const USER_TOKEN =`${USER_PREFIX}token`;
const USER_INFO=`${USER_PREFIX}info`;

//以下两种写法等价
// function set(key,data) {
// }
const set = (key,data) =>{
    localStorage.setItem(key,data);
};
const get = (key) => {
    return localStorage.getItem(key);
};
//对外暴露
export default {
    set,
    get,
    USER_INFO,
    USER_TOKEN,
};

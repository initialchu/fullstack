//此文件是认证相关逻辑
import {defineStore} from 'pinia'

import axios from '../axios'
import {ref} from 'vue'
import {useRouter} from 'vue-router'

 const router = useRouter()
export const useAuthStore = defineStore('auth',()=>{
    const token = ref<string|null>(localStorage.getItem('token'))
    // 登录方法
    const login = async(username:string,password:string)=>{
        try{
            if(!username||!password){
                alert('请输入用户名和密码')
                return
            }
            // 发送登录请求
            const res=await axios.post('/auth/login',{username,password})
            // 登录成功后，获取token
            token.value = res.data.token
            // 将token存储在localStorage中
            localStorage.setItem('token',token.value||'')
            // 跳转到首页（由调用者处理）
               router.push({name:'home'})
        }catch(err:any){
            const msg = err.response?.data?.error 
        alert(msg)
            // console.error('登录失败:',  ElMessage.error('登录失败，请检查用户名和密码'),error)
        }
    }
// 注册方法
    const register = async(username:string,password:string)=>{
        try{
            if(!username||!password){
                alert('请输入用户名和密码')
                return
            }
            // 发送注册请求
            const res=await axios.post('/auth/register',{username,password})
            // 注册成功后，获取token
            token.value = res.data.token
            // 将token存储在localStorage中
            localStorage.setItem('token',token.value||'')
            // 跳转到首页（由调用者处理）
             router.push({name:'home'})
        }catch(err:any){
           const msg = err.response?.data?.error
           alert(msg)
        }
    }

// 退出登录方法
const logout = ()=>{

    token.value = null// 清除token
    localStorage.removeItem('token')// 从localStorage中移除token

}
return { token, login, register, logout }
})
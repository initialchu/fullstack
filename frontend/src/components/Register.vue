<template>
    <div class="login-container">
        <el-container>
            <el-main>
                <div class="auth-login">
                <el-form :model="form" class="auth-login-form">
                    <h2>注册</h2>
                    <el-form-item label="用户名" label-width="80px">
                        <el-input v-model="form.username" placeholder="请输入用户名"></el-input>

                    </el-form-item>
                    <el-form-item label="密码" label-width="80px">
                        <el-input v-model="form.password" type="password" placeholder="请输入密码"></el-input>
                    </el-form-item>
                    <el-form-item  class="submit-btn">
                        <el-button @click="register"  type="primary">登录</el-button>
                    </el-form-item>
                </el-form>
                </div>
            </el-main>
        </el-container>
    </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import axios from '../axios'
import {useRouter} from 'vue-router'

const router = useRouter()
const form = ref({
    username:'',
    password:''
})
const token = ref<string|null>(localStorage.getItem('token'))

const register = async()=>{
    const username = form.value.username
    const password = form.value.password
   if(!username||!password){
    alert('请输入用户名和密码')
    return
   }
   // 发送登录请求
   const res=await axios.post('/auth/register',{username,password})
   // 注册成功后，获取token
   token.value = res.data.token
   // 将token存储在localStorage中
   localStorage.setItem('token',token.value||'')
   // 跳转到首页
   
   router.push({name:'home'})

}
</script>

<style scoped>
.auth-login-form .el-input{
    height:40px;
    width:70%;
    font-size:18px;
    
}
.login-container{
    height:100%;
}
.auth-login{
    text-align:center;
    display:flex;
    justify-content:center;
    align-items:center;
    
}
.auth-login-form{
    font-size:25px;
    min-width:600px;
    height:400px;
   max-width:800px;
    padding:20px;
    border:1px solid #ccc;
    border-radius:5px;
   
    
}
.auth-login-form .el-form-item{
    margin-bottom:50px;
    

}
.submit-btn{
   width:100%;
   
   display:flex;
   flex-direction:column;
    align-items:center;
   

}
.submit-btn .el-button{
    width:200px;
    height:40px;
    font-size:18px;
    
    letter-spacing:20px;
}

   
    

</style>
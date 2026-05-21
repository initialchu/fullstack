<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <el-menu
          :default-active="activeIndex"
          class="el-menu-demo"
          mode="horizontal"
          :ellipsis="false"
          @select="handleSelect"
        >
          
          <el-menu-item index="home">首页</el-menu-item>
          <el-menu-item index="currencyexchange">兑换货币</el-menu-item>
          <el-menu-item index="news">查看新闻</el-menu-item>
          <el-menu-item index="login" v-if="!authStore.isLoggedIn">登录</el-menu-item>
          <el-menu-item index="register" v-if="!authStore.isLoggedIn">注册</el-menu-item>
          <el-menu-item index="logout" v-if="authStore.isLoggedIn">退出</el-menu-item>
          
           
          
        </el-menu>
      </el-header>
      <el-main>
       <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref ,watch} from 'vue'
import {useRouter,useRoute} from 'vue-router'
import {useAuthStore} from './store/auth'
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const activeIndex = ref('home')
const handleSelect = (key: string) => {
  if(key ==='logout'){
    authStore.logout()
    router.push({name:'home'})
    return
  }else{
  console.log(key)
  router.push({
    name:key
  })
}
}
watch(route,(newRoute)=>{
  if(newRoute.name){
    activeIndex.value = newRoute.name.toString()
  }
  activeIndex.value = "home"
})

</script>

<style scoped>
.el-menu--horizontal > .el-menu-item {
 font-size:18px;
}
</style>
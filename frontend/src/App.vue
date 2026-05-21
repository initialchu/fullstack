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
          <el-menu-item index="login">登录</el-menu-item>
          <el-menu-item index="register">注册</el-menu-item>
          <el-menu-item index="logout">退出</el-menu-item>
          
           
          
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
const router = useRouter()
const route = useRoute()

const activeIndex = ref('home')
const handleSelect = (key: string) => {
  console.log(key)
  router.push({
    name:key
  })
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
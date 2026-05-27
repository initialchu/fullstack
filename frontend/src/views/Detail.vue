<template>
    <el-container >
        <el-main>
             
  <div >
   
    <el-card class="article-card">
       <h2>{{article.Title}}</h2>
        <article>{{article.Content}}</article>
        <el-button @click="router.back()">返回</el-button>
         <footer >&copy; 2026 My  App</footer>
    </el-card>
   
   
  </div>
 
  </el-main>
  </el-container>
</template>

<script setup lang="ts">
import {useRoute} from 'vue-router'
import {useRouter} from 'vue-router'
import {ref} from 'vue'
import axios from '../axios'
import {onMounted} from 'vue'
const router = useRouter()
const route = useRoute()
const article = ref({
    ID:'',
    Title:'',
    Content:'',
    
})
const getDetail = async ()=>{
    try{
        const res = await axios.get(`/articles/${route.query.id}`)
        article.value = res.data
        console.log('article',res.data)

    }catch(err:any){
        const msg=err.response?.data?.err
        alert(msg||'获取文章详情失败')
    }
}
onMounted(()=>{
    getDetail()
})

</script>

<style scoped>
div{
    display:flex;
    justify-content:center;
    padding:20px;
}
.article-card{
    max-width:1000px;
    min-width:600px;
    min-height:600px;
    font-size:18px;
    margin-bottom:20px;
    letter-spacing:1.5px;
    background-color:#fff;
    position:relative;
}
footer{
    margin-top:20px;
    padding:10px;
    position:absolute;
    bottom:0;
    right:0;
    color:#888;
}
.el-button{
    position:absolute;
    margin:20px;
    right:20px;
    font-size:16px;
    border:1px solid #000000;
}

</style>
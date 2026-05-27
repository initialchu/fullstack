<template>
    <el-container >
        <el-main>
  <div v-if="articles&&articles.length">
    
    <el-card @click="goDetail(article.ID)" class="article-card" v-for="article in articles" :key="article.ID">
        <h2>{{ article.Title }}</h2>
        <p>{{ article.Preview }}</p>
        <el-button text @click="goDetail(article.ID)">阅读更多</el-button>
        <el-button class="likes" @click.stop="like(article.ID)"></el-button>
        <span>{{ likes.get(article.ID) }}</span>
    </el-card>
   
  </div>
  <div class="no-data" v-else>no-data</div>
  </el-main>
  </el-container>
</template>

<script setup lang="ts">
import axios from '../axios'
import{ref, onMounted} from 'vue'
import {useRouter} from 'vue-router'

const router = useRouter()
interface Article {
    ID:string;
    Title:string;
    Content:string;
    Preview:string;
}
const articles = ref<Article[]>([])
const goDetail = (id:string)=>{
    router.push({name:'detail',query:{id}})
}
// 点赞
const like = async (id:string)=>{
    try{
        await axios.post(`articles/${id}/like`)
    }catch(err:any){
        const msg=err.response?.data?.err
        alert(msg||'点赞失败')
    }
    fetchlikes(id)
}

// 获取点赞
const likes = ref(new Map<string, number>())
 const fetchlikes = async (id:string)=>{
    try{
        const res = await axios.get(`articles/${id}/like`)
        likes.value.set(id,res.data.likes)
        
         console.log('likes',res.data)
    }catch(err:any){
        const msg=err.response?.data?.err
        alert(msg||'获取点赞数失败')
    }
   
 }
const fetchArticles = async ()=>{
    const res = await axios.get<Article[]>('/articles')
    articles.value = res.data
    articles.value.forEach(article => {
        fetchlikes(article.ID)
    })
    console.log('articles',articles.value)
}
 onMounted(fetchArticles)
</script>

<style scoped>
.likes{
    background:url('../assets/like.png') no-repeat center;
    border-radius:50%;
    background-size:cover;
    
}
.article-card{
    margin-bottom:20px;
}
</style>
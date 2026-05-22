<template>
    <el-container >
        <el-main>
  <div v-if="articles&&articles.length">
    
    <el-card @click="goDetail(article.ID, article.Title, article.Content)" class="article-card" v-for="article in articles" :key="article.ID">
        <h2>{{ article.Title }}</h2>
        <p>{{ article.Preview }}</p>
        <el-button text @click="goDetail(article.ID, article.Title, article.Content)">阅读更多</el-button>
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
const goDetail = (id:string,title:string,content:string)=>{
    router.push({name:'detail',query:{id,title,content}})
}

const fetchArticles = async ()=>{
    const res = await axios.get<Article[]>('/articles')
    articles.value = res.data
    
    console.log('articles',articles.value)
}
 onMounted(fetchArticles)
</script>

<style scoped>
.article-card{
    margin-bottom:20px;
}
</style>
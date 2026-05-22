import {createRouter, createWebHistory} from 'vue-router'
import {type RouteRecordRaw} from 'vue-router'

const routes:RouteRecordRaw[] = [
    {
        path:'/',
        name:'home',
        component:()=>import('../views/Home.vue')
    },
    {
        path:'/currencyexchange',
        name:'currencyexchange',
        component:()=>import('../views/CurrencyExchange.vue')
    },
    {
        path:'/news',
        name:'news',
        component:()=>import('../views/News.vue')

    },
    {
        path:'/login',
        name:'login',
        component:()=>import('../components/Login.vue')
    },
    {
        path:'/register',
        name:'register',
        component:()=>import('../components/Register.vue')
    },
    {
        path:'/detail',
        name:'detail',
        component:()=>import('../views/Detail.vue')
    }
]



const router = createRouter({
    history:createWebHistory(),
    routes,
})

export default router
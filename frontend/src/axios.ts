import axios from 'axios';

const instance = axios.create({
    baseURL:'http://localhost:3000/api',
})

instance.interceptors.request.use(config =>{
    const token = localStorage.getItem('token');
    if(token){
        // ensure headers object exists and set Authorization with a space after Bearer
        config.headers = config.headers || {};
        (config.headers as any).Authorization = 'Bearer ' + token;
    }
    return config;
})

export default instance;
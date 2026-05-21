<template>
<el-container class="exchange-container">
    <div class="Form">
    <el-header>
        <h2>货币兑换</h2>
    </el-header>
    <el-main>
        <exchange-form></exchange-form>
    </el-main>
  <el-form :model="form" class="exchange-form">
    <el-form-item label="从货币" label-width="80px">
       <el-select v-model="form.from" placeholder="请选择货币">
        <el-option v-for="cu in currencies" :label="cu" :value="cu" />
       
        </el-select>
    </el-form-item>
    <el-form-item label="到货币" label-width="80px">
      <el-select v-model="form.to" placeholder="请选择货币">
        <el-option v-for="cu in currencies" :label="cu" :value="cu" />
        </el-select>
      
    </el-form-item>
    <el-form-item label="金额" label-width="80px">
     
      <el-input v-model="form.amount" type="number" palceholder="请输入金额"></el-input>
    </el-form-item>
    <el-form-item >
        <el-button @click=" onSubmit" class="submit-btn" type="primary">兑换</el-button>
    </el-form-item>

  </el-form>
  <div v-if="result" class="result">
    <h3>兑换结果：{{ result }}</h3>
    
  </div>
  </div>
</el-container>
</template>

<script lang="ts" setup>
import { reactive,ref ,onMounted} from 'vue'
import axios from '../axios'


interface ExchangeInfo {
   fromCurrency: string; 
   toCurrency: string;
    rate: number;
}
const rates = ref<ExchangeInfo[]>([])//汇率数据
const currencies = ref<string[]>([])//货币列表
const result = ref<number>(0)//兑换结果
const fetchCurrencies = async ()=>{
const res=await axios.get<ExchangeInfo[]>('/exchangerate	')
rates.value = res.data
console.log('ooo',rates.value)
currencies.value=[...new Set(res.data.map((rate: ExchangeInfo) => [rate.fromCurrency, rate.toCurrency]).flat())]
}

onMounted(fetchCurrencies)

const form = reactive({
    from:'',
    to:'',
 
    amount:0,
  
})

const onSubmit = () => {
  const rate = rates.value.find((rate)=>rate.fromCurrency===form.from && rate.toCurrency === form.to)
    if(rate){
        result.value=Number((rate.rate*form.amount).toFixed(2))
    }else{
        result.value=0
    }
  
}
</script>

<style scoped>
.exchange-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  
  
}
.Form {
  min-width: 600px;
  max-width: 800px;
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
}
.exchange-form {
  margin-top: 20px;
  position:relative;
}
.submit-btn {
    max-width: 200px;
  text-align: center;
  position:relative;
  left:45%;
}
</style>
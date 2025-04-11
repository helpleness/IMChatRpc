import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import Primevue from 'Primevue/config'
import Aura from '@primevue/themes/aura'
import Button from "primevue/button";

import Card from 'primevue/card';

import {definePreset} from "@primevue/themes";

import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Checkbox from 'primevue/checkbox';
import Password from 'primevue/password';
import SpeedDial from 'primevue/speeddial';
import ScrollPanel from 'primevue/scrollpanel';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';

import FileUpload from 'primevue/fileupload';
import ToastService from 'primevue/toastservice';  // 添加这个导入
import Toast from 'primevue/toast';  // 添加这个导入


import Avatar  from 'primevue/avatar';
import 'primeicons/primeicons.css'

const app = createApp(App)


const MyPreset = definePreset(Aura, {
    semantic: {
        primary: {
            50: '{red.50}',
            100: '{red.100}',
            200: '{red.200}',
            300: '{red.300}',
            400: '{red.400}',
            500: '{red.500}',
            600: '{red.600}',
            700: '{red.700}',
            800: '{red.800}',
            900: '{red.900}',
            950: '{red.950}'
        }
    }
});

app.use(createPinia())
app.use(ToastService);
app.use(router)
app.use(Primevue,{
    theme: {
        preset: MyPreset
    }
});
app.component('Button',Button)
app.component("Card",Card)
app.component("InputGroup",InputGroup)
app.component("InputGroupAddon",InputGroupAddon)
app.component("InputText",InputText)
app.component("InputNumber",InputNumber)
app.component("Password",Password)
app.component("Checkbox",Checkbox)
app.component("ScrollPanel",ScrollPanel)
app.component("Splitter",Splitter)
app.component("SplitterPanel",SplitterPanel)
app.component("FileUpload",FileUpload)
app.component('Toast', Toast);  // 注册 Toast 组件
app.component("Avatar",Avatar)
app.component("SpeedDial",SpeedDial)
app.mount('#app')

import Vue from 'vue'
import VueResource from 'vue-resource'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-default/index.css'
import App from './App.vue'

Vue.use(VueResource);
Vue.use(ElementUI);

Vue.config.productionTip = false;
Vue.http.options.xhr = {withCredentials: true};
Vue.http.options.root = 'https://127.0.0.1:3333'; // TODO

new Vue({
    el: '#app',
    render: h => h(App)
});

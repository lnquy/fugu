import Vue from 'vue'
import VueResource from 'vue-resource'
import ElementUI from 'element-ui'
import VueCodeMirror from 'vue-codemirror'
import 'element-ui/lib/theme-default/index.css'
import App from './App.vue'
import {store} from './store/store.js'

Vue.use(VueResource);
Vue.use(ElementUI);
Vue.use(VueCodeMirror);

Vue.config.productionTip = false;
Vue.http.options.xhr = {withCredentials: true};
Vue.http.options.root = 'https://127.0.0.1:3333'; // TODO

new Vue({
    el: '#app',
    store: store,
    render: h => h(App)
});

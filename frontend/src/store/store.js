import Vue from 'vue';
import Vuex from 'vuex';
Vue.use(Vuex);

import Architecture from './modules/architecture';
import Language from './modules/language';
import OptimizeDialog from './modules/optimize_dialog';

export const store = new Vuex.Store({
    modules: {
        Architecture,
        Language,
        OptimizeDialog,
    }
});

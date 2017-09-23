const state = {
    optmd_show: false,
    optmd_data: {
        name: '',
            fields: [],
            info: {
            text: ''
        }
    },
};

const getters = {
    optmd_show: (state) => {
        return state.optmd_show;
    },
    optmd_data: (state) => {
        return state.optmd_data;
    }
};

const mutations = {
    setOptmdShow: (state, payload) => {
        state.optmd_show = payload;
    },
    setOptmdData: (state, payload) => {
        state.optmd_data = payload;
    },
};

export default {
    state,
    getters,
    mutations
}



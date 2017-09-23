const state = {
    lang: 'go',
};

const getters = {
    lang: (state) => {
        return state.lang;
    },
};

const mutations = {
    setLang: (state, payload) => {
        state.lang = payload;
    },
};

export default {
    state,
    getters,
    mutations
}



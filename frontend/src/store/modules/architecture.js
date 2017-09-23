const state = {
    arch: 'amd64',
};

const getters = {
    arch: (state) => {
        return state.arch;
    },
};

const mutations = {
    setArch: (state, payload) => {
        state.arch = payload;
    },
};

export default {
    state,
    getters,
    mutations
}



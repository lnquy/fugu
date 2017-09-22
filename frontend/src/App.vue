<template>
    <div id="app">
        <el-menu theme="dark" :default-active="activeIndex" class="el-menu-demo" mode="horizontal" @select="handleSelect">
            <el-menu-item index="1" disabled="true"><span style="font-family: 'Satisfy', cursive; font-size: 32px">Fugu</span></el-menu-item>
            <el-submenu index="2">
                <template slot="title">Languages</template>
                <el-menu-item index="2-1">Go</el-menu-item>
                <el-menu-item index="2-2">C/C++</el-menu-item>
                <el-menu-item index="2-3">Java</el-menu-item>
            </el-submenu>
            <el-menu-item index="3">Architectures</el-menu-item>
        </el-menu>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                activeIndex: '1',
            }
        },
        methods: {
            handleSelect(key, keyPath) {
                console.log(key, keyPath);
            }
        }
    }
</script>

<style>
    body {
        margin: 0;
    }
    #app {
        font-family: Helvetica, sans-serif;
        text-align: center;
    }
    .el-menu {
        border-radius: 0;
        text-align: left;
    }
</style>

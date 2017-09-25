<template>
    <el-tooltip class="item" effect="dark" content="Padding byte" placement="top">
        <div class="fg-box fg-padding-box"></div>
    </el-tooltip>
</template>

<script>
    export default {
        data() {
            return {}
        }
    }
</script>


<style scoped>
    .fg-box {
        width: 15px;
        height: 15px;
        margin: 5px;
    }

    .fg-padding-box {
        background-color: #F44336;
        display: inline-flex;
        flex-wrap: wrap;
    }
</style>

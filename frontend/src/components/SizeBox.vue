<template>
    <el-tooltip class="item" effect="dark" content="Size byte" placement="top">
        <div class="fg-box fg-size-box"></div>
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

    .fg-size-box {
        background-color: #4CAF50;
        display: inline-flex;
        flex-wrap: wrap;
    }
</style>

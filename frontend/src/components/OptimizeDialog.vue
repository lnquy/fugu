<template>
    <el-dialog :title="optmd_data.name" :visible.sync="optmd_show" size="tiny">
        <textarea readonly>{{ optmd_data }}</textarea>
        <span slot="footer" class="dialog-footer">
            <el-button type="primary" @click="optmd_show = false">Ok</el-button>
        </span>
    </el-dialog>
</template>

<script>
    export default {
        data() {
            return {
            };
        },
        props: ["optmd_show", "optmd_data"],
    };
</script>

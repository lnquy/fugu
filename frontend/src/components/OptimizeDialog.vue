<template>
    <el-dialog :title="optmd_data.name | title" :visible.sync="optmd_show" :before-close="handleClose" :close-on-click-modal="false">
        <pre class="fg-pre">{{ optmd_data.info.text }}</pre>
        <hr class="fg-hr">
        <result-table :tbl="optmd_data"></result-table>
        <span slot="footer" class="dialog-footer">
            <el-button type="primary" @click="handleClose">Ok</el-button>
        </span>
    </el-dialog>
</template>

<script>
    import ResultTable from "./ResultTable.vue";
    import {mapGetters} from "vuex";
    import {mapMutations} from "vuex";

    export default {
        data() {
            return {
                tmp: {
                    name: '',
                    fields: [],
                    info: {
                        text: ''
                    }
                }
            };
        },
        components: {
            ResultTable,
            'result-table': ResultTable
        },
        computed: {
            ...mapGetters([
                'optmd_show',
                'optmd_data',
            ])
        },
        methods: {
            ...mapMutations([
                'setOptmdShow',
                'setOptmdData',
            ]),
            handleClose() {
                this.setOptmdShow(false);
                this.setOptmdData(this.tmp);
            }
        },
        filters: {
            title(val) {
                return 'Optimized ' + val + ' struct';
            }
        }
    }
</script>

<style>
    .el-dialog__body {
        padding: 20px 20px 0 20px;
    }
</style>

<style scoped>
    .fg-pre {
        padding-bottom: 10px;
    }

    .fg-hr {
        border: 0;
        height: 1px;
        background-image: linear-gradient(to right, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.75), rgba(0, 0, 0, 0));
    }

</style>

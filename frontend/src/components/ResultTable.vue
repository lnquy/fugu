<template>
    <div>
        <div style="padding: 16px 0;">
            <span style="font-size: 24px">
                {{ tbl.name }}
                <span v-if="tbl.info.optimizable == false" style="color: #4CAF50; font-size: 18px">
                    <el-tooltip class="item" effect="dark" content="Struct memory is optimized" placement="top">
                        <i class="el-icon-circle-check"></i>
                    </el-tooltip>
                </span>
                <span v-else style="color: #F44336; font-size: 18px">
                    <el-tooltip class="item" effect="dark" placement="top">
                        <div slot="content">Too much padding bytes.<br>Click to optimize this struct's memory!</div>
                        <i class="el-icon-circle-cross ic-optimize" @click="optimizeStruct(tbl)"></i>
                    </el-tooltip>
                </span>
            </span>
            <span style="position: absolute; right: 30px; color: #aaa">
                <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">The actual allocated bytes on memory.<br>Actual size = Struct + Padding</div>
                    <span style="font-size: 12px">Actual: <span style="color: #2196F3">{{ tbl.info.total_size }}</span></span>
                </el-tooltip>&nbsp; -&nbsp;
                <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">Memory size used by struct (in theory)</div>
                    <span style="font-size: 12px">Struct: <span style="color: #4CAF50">{{ tbl.info.total_size - tbl.info.total_padding }}</span></span>
                </el-tooltip>&nbsp; -&nbsp;
                <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">Padding bytes for aligned struct fields</div>
                    <span style="font-size: 12px">Padding: <span style="color: #F44336">{{ tbl.info.total_padding }}</span></span>
                </el-tooltip>
            </span>
        </div>

        <table class="fg-res-table">
            <thead>
            <tr>
                <th class="text-center">Field</th>
                <th class="text-center">Type</th>
                <th class="text-center">Byte</th>
                <th>Memory alignment</th>
            </tr>
            </thead>
            <tr v-for="f in tbl.fields" :key="f.name">
                <td class="text-center">{{ f.name }}</td>
                <td class="text-center">{{ f.type }}</td>
                <td class="text-center">{{ f.size }}</td>
                <td style="display: flex; flex-wrap:wrap;">
                    <index-box v-for="i in f.index" :key="i"></index-box>
                    <span v-if="f.size <= getChunkByte()">
                        <size-box v-for="i in f.size" :key="i"></size-box><padding-box v-for="i in f.padding" :key="i"></padding-box>
                    </span>
                    <span v-else>
                        <span v-if="f.size> getChunkByte()*8">
                            <span style="font-size: 12px; padding-left: 5px;">First {{ omittedBytes(f.size) }} bytes omitted...<br></span>
                            <span v-for="i in 8" :key="i">
                                <size-box v-for="i in getChunkByte()" :key="i"></size-box><br/>
                            </span>
                            <size-box v-for="i in f.size%getChunkByte()" :key="i"></size-box><padding-box v-for="i in f.padding" :key="i"></padding-box>
                        </span>
                        <span v-else>
                            <span v-for="i in f.size/getChunkByte() >> 0" :key="i">
                                <size-box v-for="j in getChunkByte()" :key="j"></size-box><br/>
                            </span>
                            <size-box v-for="i in f.size%getChunkByte()" :key="i"></size-box><padding-box v-for="i in f.padding" :key="i"></padding-box>
                        </span>
                    </span>
                </td>
            </tr>
        </table>
    </div>
</template>

<script>
    import IndexBox from './IndexBox.vue';
    import SizeBox from './SizeBox.vue';
    import PaddingBox from './PaddingBox.vue';
    import {mapGetters} from "vuex";
    import {mapMutations} from "vuex";

    export default {
        data() {
            return {}
        },
        props: ["tbl"],
        components: {
            'index-box': IndexBox,
            'size-box': SizeBox,
            'padding-box': PaddingBox,
        },
        computed: {
            ...mapGetters([
                'lang',
                'arch',
            ])
        },
        methods: {
            ...mapMutations([
                'setOptmdShow',
                'setOptmdData',
            ]),
            getChunkByte() {
                if (this.arch === "i386") {
                    return 4
                }
                if (this.arch === "amd64") {
                    return 8
                }
            },
            optimizeStruct(val) {
                this.$http.post("api/v1/fugu/lang/" + this.lang + "/arch/" + this.arch + "/optimize", val).then(resp => {
                    this.setOptmdData(resp.body);
                    this.setOptmdShow(true);
                }, err => {
                    this.$notify.error({
                        title: 'Error',
                        message: "Failed to optimize struct:\n" + err.bodyText,
                        duration: 0
                    });
                    console.log(err); // TODO
                });
            },
            omittedBytes(val) {
                let chunk = this.getChunkByte();
                return ((val / chunk >> 0) - 8) * chunk;
            }
        },
    }
</script>

<style scoped>
    .ic-optimize:hover,
    .ic-optimize:focus,
    .ic-optimize:active {
        cursor: pointer;
    }

    /*** Table Styles **/

    .fg-res-table {
        background: white;
        border-radius: 4px;
        border-collapse: collapse;
        width: 100%;
        height: 100%;
        box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
        animation: float 5s infinite;
        text-align: left;
        margin-bottom: 20px;
    }

    th {
        color: #D5DDE5;;
        background: #263238;
        /*border-bottom: 4px solid #9ea7af;*/
        border-right: 1px solid #343a45;
        font-size: 16px;
        font-weight: 100;
        padding: 20px;
        text-align: left;
        text-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
        vertical-align: middle;
    }

    th:first-child {
        border-top-left-radius: 3px;
    }

    th:last-child {
        border-top-right-radius: 3px;
        border-right: none;
    }

    tr {
        border-top: 1px solid #C1C3D1;
        border-bottom: 1px solid #C1C3D1;
        color: #666B85;
        font-size: 14px;
        font-weight: normal;
        text-shadow: 0 1px 1px rgba(256, 256, 256, 0.1);
    }

    /*tr:hover td {*/
    /*background: #4E5066;*/
    /*color: #FFFFFF;*/
    /*border-top: 1px solid #22262e;*/
    /*}*/

    tr:first-child {
        border-top: none;
    }

    tr:last-child {
        border-bottom: none;
    }

    tr:nth-child(odd) td {
        background: #f6f6f6;
    }

    /*tr:nth-child(odd):hover td {*/
    /*background: #4E5066;*/
    /*}*/

    tr:last-child td:first-child {
        border-bottom-left-radius: 3px;
    }

    tr:last-child td:last-child {
        border-bottom-right-radius: 3px;
    }

    td {
        background: #FFFFFF;
        padding: 10px 20px;
        text-align: left;
        vertical-align: middle;
        font-weight: 300;
        font-size: 16px;
        text-shadow: -1px -1px 1px rgba(0, 0, 0, 0.1);
        border-right: 1px solid #C1C3D1;
    }

    td:last-child {
        border-right: 0;
    }

    th.text-left {
        text-align: left;
    }

    th.text-center {
        text-align: center;
    }

    th.text-right {
        text-align: right;
    }

    td.text-left {
        text-align: left;
    }

    td.text-center {
        text-align: center;
    }

    td.text-right {
        text-align: right;
    }

</style>

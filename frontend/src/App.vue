<template>
    <div id="app">
        <el-menu theme="dark" :default-active="activeIndex" class="el-menu-demo" mode="horizontal"
                 @select="handleSelect">
            <el-menu-item index="1"><span style="font-family: 'Satisfy', cursive; font-size: 30px;">Fugu</span>
            </el-menu-item>
            <!--<el-submenu index="2">-->
            <!--<template slot="title">Languages</template>-->
            <!--<el-menu-item index="2-1">Go</el-menu-item>-->
            <!--<el-menu-item index="2-2">C/C++</el-menu-item>-->
            <!--<el-menu-item index="2-3">Java</el-menu-item>-->
            <!--</el-submenu>-->
            <!--<el-menu-item index="3">Architectures</el-menu-item>-->
        </el-menu>
        <el-row>
            <el-col :span="12" class="fg-left-panel">
                <el-form :inline="true" :model="fuguForm" class="fg-form">
                    <el-form-item label="Language">
                        <el-select v-model="fuguForm.language" placeholder="Choose programming language">
                            <el-option label="Go" value="go"></el-option>
                            <el-option label="C/C++" value="c/c++"></el-option>
                            <el-option label="Java" value="java"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="Architecture">
                        <el-select v-model="fuguForm.arch" placeholder="Choose architecture">
                            <el-option label="i386 (32 bit)" value="i386"></el-option>
                            <el-option label="amd64 (64 bit)" value="amd64"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="onSubmit">Go</el-button>
                    </el-form-item>
                </el-form>
                <div class="codemirror">
                    <codemirror v-model="ta_code" :options="editorOptions"></codemirror>
                </div>
            </el-col>
            <transition name="el-fade-in">
                <el-col v-if="res_table.length != 0" :span="12" class="fg-right-panel">
                <div v-for="tbl in res_table" :key="tbl.name">
                    <div style="padding: 16px 0;">
                        <span style="font-size: 24px">
                            {{ tbl.name }}
                            <span v-if="tbl.info.optimizable == false" style="color: green; font-size: 18px">
                                <el-tooltip class="item" effect="dark" content="Struct memory is optimized" placement="top">
                                    <i class="el-icon-circle-check"></i>
                                </el-tooltip>
                            </span>
                            <span v-else style="color: red; font-size: 18px">
                                <el-tooltip class="item" effect="dark" placement="top">
                                    <div slot="content">Too much padding bytes.<br>Click to optimize this struct memory!</div>
                                    <i class="el-icon-circle-cross ic-optimize" @click="optimizeStruct(tbl)"></i>
                                </el-tooltip>
                            </span>
                        </span>
                        <span style="position: absolute; right: 30px; color: #aaa">
                            <span style="font-size: 12px">Actual size: <span style="color: #4CAF50">{{ tbl.info.total_size }}</span></span>&nbsp; -&nbsp;
                            <span style="font-size: 12px">Padding: <span style="color: #F44336">{{ tbl.info.total_padding }}</span></span>
                        </span>
                    </div>

                    <table class="fg-res-table">
                        <thead>
                        <tr>
                            <th class="text-center">Field</th>
                            <th class="text-center">Byte</th>
                            <th>Memory alignment</th>
                        </tr>
                        </thead>
                        <tr v-for="f in tbl.fields" :key="f.name">
                            <td class="text-center">{{ f.name }}</td>
                            <td class="text-center">{{ f.size }}</td>
                            <td style="display: flex; flex-wrap:wrap;">
                                <index-box v-for="i in f.index" :key="i"></index-box>
                                <span v-if="f.size <= getChunkByte()">
                                    <size-box v-for="i in f.size" :key="i"></size-box>
                                </span>
                                <span v-else>
                                    <span v-for="i in f.size/getChunkByte()" :key="i">
                                        <size-box v-for="i in getChunkByte()" :key="i"></size-box><br>
                                    </span>
                                    <size-box v-for="i in f.size%getChunkByte()" :key="i"></size-box>
                                </span>
                                <padding-box v-for="i in f.padding" :key="i"></padding-box>
                            </td>
                        </tr>
                    </table>
                </div>
            </el-col>
            </transition>
        </el-row>

        <optimize-dialog :optmd_show="optm_dialog.show" :optmd_data="optm_dialog.data"></optimize-dialog>
    </div>
</template>

<script>
    import IndexBox from './components/IndexBox.vue';
    import SizeBox from './components/SizeBox.vue';
    import PaddingBox from './components/PaddingBox.vue';
    import OptimizeDialog from './components/OptimizeDialog.vue';

    export default {
        data() {
            return {
                activeIndex: '1',
                fuguForm: {
                    language: 'go',
                    arch: 'amd64',
                },
                ta_code: '',
                res_table: [],
                editorOptions: {
                    tabSize: 4,
                    mode: 'text/x-go',
                    theme: 'material',
                    lineNumbers: true,
                    line: true,
                    placeholder: 'type MyStruct struct {...}',
                    foldGutter: true,
                    gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
                    styleSelectedText: true,
                    highlightSelectionMatches: {showToken: /\w/, annotateScrollbar: true},
                },
                optm_dialog: {
                    show: false,
                    data: '',
                }
            }
        },
        components: {
            'index-box': IndexBox,
            'size-box': SizeBox,
            'padding-box': PaddingBox,
            'optimize-dialog': OptimizeDialog,
        },
        methods: {
            handleSelect(key, keyPath) {
                console.log(key, keyPath);
            },
            onSubmit() {
                this.$http.post("api/v1/fugu/lang/" + this.fuguForm.language + "/arch/" + this.fuguForm.arch, this.ta_code).then(resp => {
                    if (resp.body[0] !== undefined) {
                        this.res_table = resp.body;
                    } else {
                        this.res_table = []
                    }
                }, err => {
                    console.log(err)
                });
            },
            getChunkByte() {
                if (this.fuguForm.arch === "i386") {
                    return 4
                }
                if (this.fuguForm.arch === "amd64") {
                    return 8
                }
            },
            optimizeStruct(val) {
                this.$http.post("api/v1/fugu/lang/" + this.fuguForm.language + "/arch/" + this.fuguForm.arch + "/optimize", val).then(resp => {
                    this.optm_dialog.data = resp.body;
                    this.optm_dialog.show = true;
                }, err => {
                    console.log(err)
                });
            }
        },
        watch: {
            'fuguForm.language': function () {
                this.res_table = []
            },
            'fuguForm.arch': function () {
                this.res_table = []
            },
        }
    }
</script>

<style>
    body {
        margin: 0;
        font-family: Helvetica, sans-serif;
    }

    .el-menu {
        border-radius: 0;
    }

    .cm-s-material {
        min-height: 600px;
        border-radius: 4px;
    }

</style>

<style scoped>
    .fg-left-panel,
    .fg-right-panel {
        padding: 20px;
    }

    .codemirror {
        font-size: 14px;
        line-height: 1.5em;
    }

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

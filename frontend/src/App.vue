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
                    <result-table :tbl="tbl"></result-table>
                </div>
            </el-col>
            </transition>
        </el-row>

        <optimize-dialog></optimize-dialog>
    </div>
</template>

<script>
    import OptimizeDialog from './components/OptimizeDialog.vue';
    import ResultTable from './components/ResultTable.vue';
    import {mapMutations} from 'vuex';

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
            }
        },
        components: {
            'result-table': ResultTable,
            'optimize-dialog': OptimizeDialog,
        },
        methods: {
            ...mapMutations([
                'setLang',
                'setArch',
            ]),
            handleSelect(key, keyPath) {
                console.log(key, keyPath);
            },
            onSubmit() {
                this.$http.post("api/v1/fugu/lang/" + this.fuguForm.language + "/arch/" + this.fuguForm.arch, this.ta_code).then(resp => {
                    if (resp.body != null && resp.body[0] !== undefined) {
                        this.res_table = resp.body;
                    } else {
                        this.res_table = []
                    }
                }, err => {
                    this.$notify.error({
                        title: 'Error',
                        message: "Failed to calculate struct memory:\n" + err.bodyText,
                        duration: 0
                    });
                });
            },
        },
        watch: {
            'fuguForm.language': function () {
                this.setLang(this.fuguForm.language);
                this.res_table = [];
            },
            'fuguForm.arch': function () {
                this.setArch(this.fuguForm.arch);
                this.res_table = [];
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

    span, div {
        margin: 0;
        padding: 0;
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
</style>

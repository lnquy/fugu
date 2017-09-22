<template>
    <div id="app">
        <el-menu theme="dark" :default-active="activeIndex" class="el-menu-demo" mode="horizontal"
                 @select="handleSelect">
            <el-menu-item index="1"><span style="font-family: 'Satisfy', cursive; font-size: 32px">Fugu</span>
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
            <el-col v-if="res_table.length != 0" :span="12" class="fg-right-panel">
                <div v-for="tbl in res_table">
                    <h3>{{ tbl.name }}</h3>
                    <table class="fg-res-table">
                        <thead>
                            <tr>
                                <th>Field</th>
                                <th>Size</th>
                                <th>Memory alignment</th>
                            </tr>
                        </thead>
                        <tr v-for="f in tbl.fields">
                            <td>{{ f.name }}</td>
                            <td>{{ f.size }}</td>
                            <td style="display: flex; flex-wrap:wrap;">
                                <div v-for="i in f.index" class="fg-box fg-index-box"></div>

                                <span v-if="f.size <= getChunkByte()">
                                    <div v-for="i in f.size" class="fg-box fg-size-box"></div>
                                </span>
                                <span v-else>
                                    <span v-for="i in f.size/getChunkByte()">
                                        <div v-for="i in getChunkByte()" class="fg-box fg-size-box"></div><br>
                                    </span>
                                    <div v-for="i in f.size%getChunkByte()" class="fg-box fg-size-box"></div>
                                </span>

                                <!--<div v-for="i in f.size" class="fg-box fg-size-box"></div>-->
                                <div v-for="i in f.padding" class="fg-box fg-padding-box"></div>
                            </td>
                        </tr>
                    </table>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script>
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
                    highlightSelectionMatches: { showToken: /\w/, annotateScrollbar: true },
                }
            }
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
            }
        },
        watch: {
            'fuguForm.language': function() {
                this.res_table = []
            },
            'fuguForm.arch': function() {
                this.res_table = []
            },
        }
    }
</script>

<style>
    body {
        margin: 0;
    }

    #app {
        font-family: Helvetica, sans-serif;
    }

    .el-menu {
        border-radius: 0;
    }

    .cm-s-material {
        border-radius: 4px;
    }
</style>

<style scoped>
    .fg-left-panel {
        padding: 20px;
    }

    .fg-right-panel {
        padding: 20px;
    }

    .fg-form {
    }

    .fg-res-table {
        width: 100%;
        height: 100%;
    }

    .fg-box {
        width: 15px;
        height: 15px;
        margin: 5px;
    }

    .fg-index-box {
        background-color: gray;
    }

    .fg-size-box {
        background-color: green;
        display: inline-flex;
        flex-wrap: wrap;
    }

    .fg-padding-box {
        background-color: red;
    }

    .codemirror{
        font-size: 14px;
        line-height: 1.5em;
    }

    table {
        max-width: 960px;
        margin: 10px auto;
        border-radius: 4px;
    }

    thead th {
        font-weight: 400;
        background: #8a97a0;
        color: #FFF;
    }

    tr {
        background: #f4f7f8;
        border-bottom: 1px solid #FFF;
        margin-bottom: 5px;
    }

    tr:nth-child(even) {
        background: #e8eeef;
    }

    th, td {
        text-align: left;
        padding: 20px;
        font-weight: 300;
    }

    td {
        padding: 10px 20px;
    }
</style>

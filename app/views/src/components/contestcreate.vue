<template>
  <el-container>
    <!-- <el-aside width="200px" style="border-right: 1px solid silver">
      <div style="position: absolute; top: 120px" id="side-menu">
        <div>title</div>
        <div>item</div>
        <div>item</div>
        <div>item</div>
        <div>item</div>
        <div>item</div>
        <div>item</div>
        <div>item</div>
      </div>
      <el-menu>

      </el-menu>
    </el-aside> -->
    <el-main>
      <!-- <div align="left" style="font-size: 24px; font-weight: 600; margin-bottom: 10px">创建比赛</div> -->
      <el-form
        ref='form'
        :model="form"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item label="比赛标题" prop="title">
          <el-input v-model="form.title" placeholder="输入标题" size="small"></el-input>
        </el-form-item>
        <el-form-item label="起止时间" prop="dateTimes">
          <el-date-picker
            v-model="form.dateTimes"
            type="datetimerange"
            range-separator="TO"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 100%"
            size="small"
          >
          </el-date-picker>
        </el-form-item>
        <el-form-item label="简要介绍" prop="desc">
          <el-input
            type="textarea"
            :autosize="{minRows: 2, maxRows: 4}"
            v-model="form.desc"
            auto-complete="off"
            size="small"
          ></el-input>
        </el-form-item>
        <el-form-item label="详细介绍" prop="details">
          <el-input
            type="textarea"
            :autosize="{minRows: 3, maxRows: 6}"
            v-model="form.details"
            auto-complete="off"
            size="small"
          ></el-input>
        </el-form-item>
        <el-form-item label="赛制脚本" prop="script">
          <div class="cm-container">
            <codemirror v-model="form.script" :options="cmOptions" align="left"></codemirror>
          </div>
        </el-form-item>
        <el-form-item label="播放器html" prop="playback">
          <div class="cm-container">
            <codemirror v-model="form.playback" :options="cmOptions" align="left"></codemirror>
          </div>
        </el-form-item>
        <el-form-item label="是否开放报名" prop="isRegOpen">
          <div align="left">
            <el-radio :label="true" v-model="form.isRegOpen">是</el-radio>
            <el-radio :label="false" v-model="form.isRegOpen">否</el-radio>
          </div>
        </el-form-item>
        <el-form-item label="选择管理员账号" prop="moderators">
          <div align="left">
            <el-tag
              v-for="(item, key) in form.moderators"
              :key="key" closable @close="handleTagClose(item)"
              align="center"
            >{{item.handle}}</el-tag>
            <el-popover
              placement="bottom"
              trigger="click"
              width="400"
              v-model="selOpen"
            > <div>
                <el-input v-model="selInput" @input="handleInput" placeholder="输入账号进行查询" size="small"></el-input>
                <el-table :data="selOptions" v-loading="selLoading">
                  <el-table-column label="账号" min-width="150px">
                    <template slot-scope="scope">
                      <div style="font-weight: 600">{{scope.row.handle}}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="昵称" prop="nickname" min-width="150px"></el-table-column>
                  <el-table-column label="选项">
                    <template slot-scope="scope">
                      <el-button
                        size="small"
                        type="text"
                        @click="addModerator(scope.row)"
                        :disabled="checkDisabled(scope.row)"
                      >添加</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-button type="text" @click="selOpen=false">关闭</el-button>
              </div>
              <el-button plain size="small" slot="reference" icon="el-icon-plus">添加</el-button>
            </el-popover>
          </div>
        </el-form-item>
      </el-form>
      <el-row>
        <el-col :span="12">
          <el-button size="small" type="primary" @click="submitContest" style="width: 80%">提交</el-button>
        </el-col>
        <el-col :span="12">
          <el-button size="small" @click="goBack" style="width: 80%">返回</el-button>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'contest_create',
  components: {
    codemirror
  },
  created () {
    // const listener = () => {
    //   this.listenerCount += 1
    //   const count = this.listenerCount
    //   window.setTimeout(() => {
    //     this.updateSidebar(count)
    //   }, 240)
    // }
    // window.addEventListener('scroll', listener)
    // this.menuListener = listener
    this.$store.commit('setStallFlag', true)
  },
  beforeDestroy () {
    window.removeEventListener('scroll', this.menuListener)
  },
  data () {
    const dateValidator = (rule, value, callback) => {
      const t1 = Math.round(value[1].getTime() / 1000)
      const t2 = Math.round(value[0].getTime() / 1000)
      if (t1 > t2) {
        callback()
      } else {
        callback(new Error('日期范围错误'))
      }
    }
    return {
      sideBarTop: 120,
      listenerCount: 0,
      menuListener: null,
      form: {
        title: '',
        dateTimes: [],
        desc: '',
        details: '',
        script: '',
        playback: '',
        isRegOpen: false,
        moderators: []
      },
      cmOptions: {
        lineNumbers: true,
        lineWrapping: true,
        indentUnit: 2,
        tabSize: 2,
        autoCloseBrackets: true
      },
      selInput: '',
      selOptions: [],
      selLoading: false,
      selOpen: false,
      rules: {
        title: [
          { validator: this.$functions.globalValidator, trigger: 'blur' },
          { required: true, message: '请输入比赛名称', trigger: 'blur' },
          { min: 3, max: 20, message: '名称应在3-20个字符之间', trigger: 'blur' }
        ],
        dateTimes: [
          { required: true, message: '请选择起止时间' },
          { validator: dateValidator, trigger: 'blur' }
        ],
        desc: [
          { required: true, message: '请输入简介', trigger: 'blur' },
          { max: 255, message: '输入过长', trigger: 'blur' }
        ],
        details: [
          { required: true, message: '请输入详细介绍', trigger: 'blur' },
          { max: 4095, message: '输入过长', trigger: 'blur' }
        ],
        isRegOpen: [
          { required: true, message: '请选择是否开放报名', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    // updateSidebar (val) {
    //   if (val !== this.listenerCount) {
    //     return
    //   }
    //   window.requestAnimationFrame(() => {
    //     const sl = document.documentElement.scrollTop || document.body.scrollTop
    //     const el = document.getElementById('side-menu')
    //     velocity(el, { top: sl + 120 }, { duration: 120 })
    //   })
    // },
    handleTagClose (val) {
      this.form.moderators.splice(this.form.moderators.indexOf(val), 1)
    },
    handleInput (val) {
      this.selOptions = []
      this.selLoading = true
      this.searchForUser(val)
    },
    addModerator (item) {
      this.form.moderators.push(item)
    },
    checkDisabled (item) {
      let flag = false
      this.form.moderators.forEach(member => {
        if (member.handle === item.handle) {
          flag = true
        }
      })
      return flag
    },
    searchForUser (val) {
      if (val !== '') {
        this.selLoading = true
        this.$axios.get(
          '/user_search/' + val
        ).then(res => {
          this.selOptions = res.data
          this.selLoading = false
        }).catch(err => {
          console.log(err)
          this.selOptions = []
          this.selLoading = false
        })
      } else {
        this.selOptions = []
        this.selLoading = false
      }
    },
    submitContest () {
      this.$refs.form.validate(valid => {
        if (valid) {
          const loading = this.$loading({ lock: true, text: '处理中' })
          let params = {
            title: this.form.title,
            start_time: Math.round(this.form.dateTimes[0].getTime() / 1000),
            end_time: Math.round(this.form.dateTimes[1].getTime() / 1000),
            desc: this.form.desc,
            details: this.form.details,
            is_reg_open: this.form.isRegOpen,
            script: this.form.script,
            playback: this.form.playback
          }
          if (this.form.moderators.length > 0) {
            const idList = []
            this.form.moderators.forEach(member => {
              idList.push(member.id)
            })
            params.moderators = idList.join(',')
          }
          params = this.$qs.stringify(params)
          this.$axios.post(
            '/contest/create',
            params
          ).then(res => {
            loading.close()
            this.$message.success('提交成功')
            this.$store.commit('setStallFlag', false)
            this.$router.push('/')
          }).catch(err => {
            loading.close()
            console.log(err)
            this.$message.error('提交失败，请重试')
          })
        }
      })
    },
    goBack () {
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.cm-container {
  border: 1px solid #dcdfe6;
  max-height: 600px;
  font-size: 14px;
  line-height: 18px;
}
</style>

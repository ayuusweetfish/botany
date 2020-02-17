<template>
  <el-dialog v-loading="loading" title="修改密码" :visible.sync="show">
    <el-form ref="passwords" :model="passwords" :rules="rules" label-position="right" label-width="100px">
      <el-form-item prop="old" :error="errors.old" label="旧密码">
        <el-input type="password" v-model="passwords.old" placeholder="请输入原密码">
        </el-input>
      </el-form-item>
      <el-form-item prop="new" :error="errors.new" label="新密码">
        <el-input type="password" v-model="passwords.new" placeholder="请输入新密码">
        </el-input>
      </el-form-item>
      <el-form-item prop="new2" :error="errors.new2" label="确认密码">
        <el-input type="password" v-model="passwords.new2" placeholder="请确认密码">
        </el-input>
      </el-form-item>
      <el-row>
        <el-col :span="12">
          <el-button style="width: 80%" type="primary" size="small" @click="confirmChange">提交</el-button>
        </el-col>
        <el-col :span="12">
          <el-button style="width: 80%" size="small" @click="cancelChange">取消</el-button>
        </el-col>
      </el-row>
    </el-form>

  </el-dialog>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
      default: 'passwordDialog'
    },
    visible: {
      type: Boolean,
      default: false
    }
  },
  data () {
    const validator2 = (rule, value, callback) => {
      if (value !== this.password.old) {
        callback(new Error('输入不一致'))
      } else {
        callback()
      }
    }
    return {
      loading: false,
      passwords: {
        old: '',
        new: '',
        new2: ''
      },
      errors: {
        old: '',
        new: '',
        new2: ''
      },
      rules: {
        old: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { validator: this.$functions.globalValidator, trigger: 'blur' },
          { min: 3, max: 30, message: '密码长度不符合要求', trigger: 'blur' }
        ],
        new: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { validator: this.$functions.globalValidator, trigger: 'blur' },
          { min: 3, max: 30, message: '密码长度不符合要求', trigger: 'blur' }
        ],
        new2: [
          { required: true, message: '请确认新密码', trigger: 'blur' },
          { validator: this.$functions.globalValidator, trigger: 'blur' },
          { validator: validator2, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    confirmChange () {
      this.loading = true
      const params = this.$qs.stringify({
        old: this.passwords.old,
        new: this.passwords.new
      })
      this.$axios.post(
        '/user/' + this.$store.state.handle + '/password',
        params
      ).then(res => {
        this.loading = false
        this.$message.success('修改成功')
        this.show = false
      }).catch(err => {
        this.loading = false
        if (err.response.status === 400) {
          this.errors.old = '密码错误'
          this.$message.error('密码错误')
        } else if (err.response.status === 403) {
          this.$message.error('权限错误')
        } else {
          this.$message.error('修改失败')
        }
      })
    },
    cancelChange () {
      this.show = false
    }
  },
  computed: {
    show: {
      set (val) {
        this.$refs.passwords.resetFields()
        this.$emit('setVisible', val)
      },
      get () {
        return this.visible
      }
    }
  }
}
</script>

<style>

</style>

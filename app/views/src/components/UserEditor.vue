<template>
  <v-card outlined style="border: none">
    <v-card-title>修改个人信息</v-card-title>
    <v-card-text>
      <v-form ref="form">
        <v-text-field
          label="昵称"
          :rules="[validators.required, validators.minmax]"
          v-model="nickname"
          counter
        ></v-text-field>
        <v-text-field
          label="邮箱"
          :rules="[validators.required, validators.email, v=>!errEmail||errEmail]"
          v-model="email"
          counter
        ></v-text-field>
        <v-textarea
          label="签名"
          placeholder="给自己添加一个bio"
          :rules="[validators.bio]"
          v-model="bio"
          :counter="256"
        ></v-textarea>
      </v-form>
    </v-card-text>
    <v-card-actions class="justify-end">
      <v-btn text color="primary" :disabled="loading" :loading="loading" @click="submit">提交</v-btn>
      <v-btn text @click="terminate">取消</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  props: {
    info: Object,
    value: Boolean
  },
  mounted () {
    this.setData()
  },
  data: () => ({
    handle: '',
    nickname: '',
    email: '',
    bio: '',
    validators: {
      required: v => !!v || '输入不能为空',
      minmax: v => (v.length >= 3 && v.length <= 30) || '长度不符合要求',
      email: v => (/^([a-zA-Z0-9]+[-_.]?)+@([a-zA-Z0-9]+\.)+[a-z]+$/.test(v)) || '请输入格式合法的邮箱',
      bio: v => v.length <= 256 || '输入不超过256个字符'
    },
    errEmail: '',
    loading: false
  }),
  methods: {
    setData () {
      this.handle = this.info.handle
      this.nickname = this.info.nickname
      this.email = this.info.email
      this.bio = this.info.bio
    },
    terminate () {
      this.$emit('input', false)
    },
    submit () {
      if (!this.$refs.form.validate()) {
        return
      }
      this.loading = true
      const params = this.$qs.stringify({
        nickname: this.nickname,
        email: this.email,
        bio: this.bio
      })
      this.$axios.post(
        '/user/' + this.handle + '/profile/edit',
        params
      ).then(res => {
        this.loading = false
        this.$emit('success')
        this.terminate()
      }).catch(err => {
        this.loading = false
        this.$emit('fail')
        if (err.response.data.err[0] === 1) {
          this.errEmail = '邮箱已被注册'
          this.$refs.form.validate()
        }
        this.errEmail = ''
      })
    }
  }
}
</script>

<style>

</style>

<template>
  <v-card outlined style="border: none">
    <v-card-title>修改密码</v-card-title>
    <v-card-text>
      <v-form ref="form">
        <v-text-field
          v-model="password0"
          label="密码"
          :rules="[validators.required, validators.minmax, v=>!errPassword0||errPassword0]"
          :append-icon="showPassword0? 'mdi-eye' : 'mdi-eye-off'"
          @click:append="showPassword0=!showPassword0"
          :type="showPassword? 'text' : 'password'"
          counter
        ></v-text-field>
        <v-text-field
          v-model="password"
          label="新密码"
          :rules="[validators.required, validators.minmax, considerP2]"
          :append-icon="showPassword? 'mdi-eye' : 'mdi-eye-off'"
          @click:append="showPassword=!showPassword"
          :type="showPassword? 'text' : 'password'"
          counter
        ></v-text-field>
        <v-text-field
          ref="password2"
          v-model="password2"
          label="确认密码"
          :rules="[validators.required, v=>v===password||'密码输入不一致']"
          :append-icon="showPassword2? 'mdi-eye' : 'mdi-eye-off'"
          @click:append="showPassword2=!showPassword2"
          :type="showPassword? 'text' : 'password'"
          counter
        ></v-text-field>
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
    handle: String
  },
  data: el => ({
    loading: false,
    password0: '',
    password: '',
    password2: '',
    showPassword0: '',
    showPassword: '',
    showPassword2: '',
    validators: {
      required: v => !!v || '输入不能为空',
      minmax: v => (v.length >= 3 && v.length <= 30) || '长度不符合要求'
    },
    considerP2: (v => {
      if (el.password2) {
        el.$refs.password2.validate()
      }
      return true
    }) || '',
    errPassword0: ''
  }),
  methods: {
    submit () {
      if (!this.$refs.form.validate()) {
        return
      }
      this.loading = true
      const params = this.$qs.stringify({
        old: this.password0,
        new: this.password
      })
      this.$axios.post(
        '/user/' + this.handle + '/password',
        params
      ).then(res => {
        this.loading = false
        this.$emit('success')
        this.terminate()
      }).catch(err => {
        this.loading = false
        this.$emit('fail')
        if (err.response.status === 400) {
          this.errPassword0 = '密码错误'
          this.$refs.form.validate()
        }
        this.errPassword0 = ''
      })
    },
    terminate () {
      this.$emit('input', false)
    }
  }
}
</script>

<style>

</style>

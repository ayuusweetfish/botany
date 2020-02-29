<template>
  <v-menu
    v-model="menu"
    :close-on-content-click="false"
    transition="scale-transition"
    offset-y
    min-width="280px"
  >
    <template v-slot:activator="{ on }">
      <v-text-field
        v-model="full"
        :label="label"
        :prepend-icon="prependIcon"
        readonly
        :hint="hint"
        v-on="on"
        :error-messages="errMessages"
      ></v-text-field>
    </template>
    <v-date-picker
      v-model="date"
      no-title
      locale="zh-cn"
      first-day-of-week="1"
      :show-current="false"
    >
      <v-container>
        <v-row>
        <v-text-field
          ref="time"
          v-model="time"
          label="输入时间"
          outlined
          dense
          class="subtitle-1"
          v-mask="'##:##:##'"
          clearable
          :rules="[validator]"
        ></v-text-field>
        </v-row>
        <v-row justify="end">
          <v-btn text color="primary" @click="confirm">确定</v-btn>
          <v-btn text @click="menu=false">取消</v-btn>
        </v-row>
      </v-container>
    </v-date-picker>
  </v-menu>
</template>

<script>
import { mask } from 'vue-the-mask'
export default {
  props: {
    value: String,
    label: String,
    prependIcon: String,
    hint: String,
    required: {
      type: Boolean,
      default: false
    },
    min: {
      type: String,
      default: ''
    },
    max: {
      type: String,
      default: ''
    }
  },
  directives: {
    mask
  },
  mounted () {
    this.full = this.value
    this.translateDate()
  },
  watch: {
    value: function (newval, oldval) {
      this.full = newval
      this.translateDate()
    },
    menu: function (newval, oldval) {
      this.full = this.value
      this.translateDate()
      if (!newval) {
        this.validate()
        this.$emit('change')
      }
    },
    min: function () {
      this.validate()
    },
    max: function () {
      this.validate()
    }
  },
  data: () => ({
    error: false,
    errMessages: [],
    menu: false,
    full: '',
    date: '',
    time: '',
    validator: v => {
      if (!v) {
        return '请输入日期'
      }
      if (/\d\d:\d\d:\d\d/.test(v)) {
        const strs = Array.from(v.split(':'), item => {
          return parseInt(item)
        })
        if (strs[0] < 24 && strs[1] < 60 && strs[2] < 60) {
          return true
        }
      }
      return '日期格式错误'
    }
  }),
  methods: {
    translateDate () {
      if (!this.full) {
        this.date = this.$functions.dateString(new Date() / 1000)
        this.time = '00:00:00'
        if (this.$refs.time) {
          this.$refs.time.validate()
        }
        return
      }
      const split = this.full.split(' ')
      this.date = split[0]
      this.time = split[1]
      if (this.$refs.time) {
        this.$refs.time.validate()
      }
    },
    confirm () {
      if (!this.$refs.time.validate()) {
        return
      }
      this.full = this.date + ' ' + this.time
      this.$emit('input', this.full)
      this.menu = false
    },
    validate () {
      if (this.required) {
        if (this.full.length === 0) {
          this.errMessages = ['请输入日期']
        } else if (this.min.length !== 0 || this.max.length !== 0) {
          const date = new Date(this.full)
          const b1 = this.min.length === 0 || date > new Date(this.min)
          const b2 = this.max.length === 0 || date < new Date(this.max)
          if (b1 && b2) {
            this.errMessages = []
          } else {
            this.errMessages = ['日期不在有效时间段内']
          }
        } else {
          this.errMessages = []
        }
      }
      return this.errMessages.length === 0
    }
  }
}
</script>

<style>

</style>

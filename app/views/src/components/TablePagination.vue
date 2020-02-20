<template>
  <div class="d-flex flex-wrap justify-end align-start">
    <div :style="$vuetify.breakpoint.smAndDown? 'width: 360px': 'width: '+(55*maxDisplay+50)+'px'">
      <v-pagination
        v-model="page"
        :total-visible="$vuetify.breakpoint.smAndDown? 5 : maxDisplay"
        :length="getLength()"
      ></v-pagination>
    </div>
    <div style="width: 120px; margin-top: 2px" class="ml-4 mr-2">
      <v-text-field
        outlined
        single-line
        :value="page"
        dense
        class="pa-0"
        :suffix="`/ ${getLength()}`"
        @change="textInput"
      ></v-text-field>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    value: Number,
    total: Number,
    maxDisplay: {
      type: Number,
      default: 8
    },
    count: {
      type: Number,
      default: 10
    }
  },
  watch: {
    value: function (newval, oldval) {
      this.page = newval
    },
    page: function (newval, oldval) {
      this.$emit('input', newval)
    }
  },
  mounted () {
    this.page = this.value
  },
  data: () => ({
    page: 0
  }),
  methods: {
    getLength () {
      return Math.ceil((this.total - 1) / this.count) || 1
    },
    textInput (val) {
      const page = Math.abs(Math.round(val)) || 1
      const length = this.getLength()
      if (page > length) {
        this.page = length
      }
    }
  }
}
</script>

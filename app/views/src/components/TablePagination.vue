<template>
  <div>
    <div v-if="$vuetify.breakpoint.mdAndUp" class="d-flex align-start justify-end">
      <v-pagination
        v-model="page"
        :length="length"
        :total-visible="$vuetify.breakpoint.smAndDown? 5 : maxDisplay"
        class="justify-end"
        :disabled="disabled"
      ></v-pagination>
      <div style="width: 120px; margin-top: 2px" class="ml-2 mr-2">
        <v-text-field
          outlined
          single-line
          :value="page"
          dense
          class="pa-0"
          :suffix="`/ ${length}`"
          @change="textInput"
          :disabled="disabled"
        ></v-text-field>
      </div>
    </div>
    <div v-else class="d-flex justify-center">
      <v-btn text @click="previousPage" :disabled="page<=1||disabled" large color="primary">
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>
      <div style="width: 120px; margin-top: 2px">
        <v-text-field
          outlined
          single-line
          :value="page"
          dense
          class="pa-0"
          :suffix="`/ ${length}`"
          @change="textInput"
          :disabled="disabled"
        ></v-text-field>
      </div>
      <v-btn text @click="nextPage" :disabled="page>=length||disabled" large color="primary">
        <v-icon>mdi-chevron-right</v-icon>
      </v-btn>
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
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  watch: {
    value: function (newval, oldval) {
      this.page = newval
    },
    page: function (newval, oldval) {
      this.$emit('input', newval)
    },
    total: function (newval, oldval) {
      this.getLength()
    }
  },
  mounted () {
    this.getLength()
  },
  data: () => ({
    page: 0,
    length: 0
  }),
  methods: {
    getLength () {
      this.page = this.value
      this.length = Math.ceil((this.total - 1) / this.count) || 1
    },
    textInput (val) {
      const page = Math.abs(Math.round(val)) || 1
      if (page > this.length) {
        this.page = this.length
      } else {
        this.page = page
      }
    },
    nextPage () {
      if (this.page < this.length) {
        this.page += 1
      }
    },
    previousPage () {
      if (this.page > 1) {
        this.page -= 1
      }
    }
  }
}
</script>

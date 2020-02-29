<template>
  <v-overlay :value="value" :absolute="absolute" :opacity="0.4">
    <v-progress-circular :color="color" :indeterminate="indeterminate" :size="size"></v-progress-circular>
  </v-overlay>
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false
    },
    absolute: {
      type: Boolean,
      default: true
    },
    color: {
      type: String,
      default: 'primary'
    },
    indeterminate: {
      type: Boolean,
      default: true
    },
    size: {
      type: Number,
      default: 32
    }
  }
}
</script>

<style>

</style>

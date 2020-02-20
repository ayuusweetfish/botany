<template>
  <div class="d-flex justify-start align-center">
    <div>
      <v-avatar :size="avatarSize()" tile>
        <v-img :src="$axios.defaults.baseURL + '/user/' + user.handle + '/avatar'"></v-img>
      </v-avatar>
    </div>
    <div>
      <div :class="nicknameClass()">{{user.nickname}}</div>
      <div :class="handleClass()">@{{user.handle}}</div>
    </div>
    <div v-if="user.handle===$store.state.handle&&identify" class="ml-2">(我)</div>
  </div>
</template>

<script>
export default {
  props: {
    user: Object,
    size: {
      type: String,
      default: 'middle'
    },
    identify: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    avatarSize () {
      if (this.size === 'large') {
        return 56
      } else if (this.size === 'small') {
        return 40
      } else {
        return 48
      }
    },
    nicknameClass () {
      if (this.size === 'large') {
        return 'ml-4 title font-weight-bold'
      } else if (this.size === 'small') {
        return 'ml-1 body-2 font-weight-bold'
      } else {
        return 'ml-2 font-weight-bold'
      }
    },
    handleClass () {
      if (this.size === 'large') {
        return 'ml-4 subtitle'
      } else if (this.size === 'small') {
        return 'ml-1 body-2'
      } else {
        return 'ml-2'
      }
    }
  }
}
</script>

<style>

</style>

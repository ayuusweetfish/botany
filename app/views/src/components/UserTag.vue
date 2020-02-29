<template>
  <div>
    <div v-if="disabled"
      class="d-flex justify-start align-center"
    >
      <div>
        <v-avatar :size="avatarSize()" tile>
          <v-img
            :style="'border-radius:' + borderRadius() + 'px'"
            :src="src"
          ></v-img>
        </v-avatar>
      </div>
      <div>
        <div :class="nicknameClass()">{{user.nickname}}</div>
        <div :class="handleClass()">@{{user.handle}}</div>
      </div>
      <div v-if="user.handle===$store.state.handle&&identify" class="ml-2">(我)</div>
    </div>
    <router-link v-else
      class="d-flex justify-start align-center"
      style="cursor: pointer; text-decoration: none"
      :to="`/profile/${user.handle}`"
    >
      <div>
        <v-avatar :size="avatarSize()" tile>
          <v-img
            :style="'border-radius:' + borderRadius() + 'px'"
            :src="src"
          ></v-img>
        </v-avatar>
      </div>
      <div>
        <div :class="nicknameClass()">{{user.nickname}}</div>
        <div :class="handleClass()">@{{user.handle}}</div>
      </div>
      <div v-if="user.handle===$store.state.handle&&identify" class="ml-2">(我)</div>
    </router-link>
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
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  watch: {
    user: function () {
      this.getSrc()
    }
  },
  mounted () {
    this.getSrc()
  },
  methods: {
    getSrc () {
      if (this.user.handle) {
        this.src = this.$axios.defaults.baseURL + '/user/' + this.user.handle + '/avatar'
      } else {
        this.src = ''
      }
    },
    avatarSize () {
      if (this.size === 'large') {
        return 56
      } else if (this.size === 'small') {
        return 40
      } else {
        return 48
      }
    },
    borderRadius () {
      if (this.size === 'large') {
        return '6'
      } else if (this.size === 'small') {
        return '4'
      } else {
        return '5'
      }
    },
    nicknameClass () {
      if (this.size === 'large') {
        return 'ml-4 title font-weight-bold secondary--text'
      } else if (this.size === 'small') {
        return 'ml-1 body-2 font-weight-bold secondary--text'
      } else {
        return 'ml-2 font-weight-bold secondary--text'
      }
    },
    handleClass () {
      if (this.size === 'large') {
        return 'ml-4 subtitle grey--text'
      } else if (this.size === 'small') {
        return 'ml-1 body-2 grey--text'
      } else {
        return 'ml-2 grey--text'
      }
    }
  },
  data: () => ({
    src: ''
  })
}
</script>

<style>

</style>

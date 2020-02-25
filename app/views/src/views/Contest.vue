<template>
  <div>
    <router-view/>
  </div>
</template>

<script>
export default {
  watch: {
    '$route.params.cid': function (newval, oldval) {
      window.location.reload()
    }
  },
  mounted () {
    this.$axios.get(
      '/contest/' + this.$route.params.cid + '/info'
    ).then(res => {
      this.$store.commit('setContest', res.data)
    }).catch(err => {
      if (err.response.status === 404) {
        // push 404
        this.$store.commit('resetContest')
      }
    })
  }
}
</script>

<style>

</style>

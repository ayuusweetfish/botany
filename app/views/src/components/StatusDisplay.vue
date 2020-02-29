<template>
  <div v-if="mode==='submission'">
    <div v-if="status===$consts.codeStat.pending"
      class="secondary--text"
    >等待中</div>
    <div v-else-if="status===$consts.codeStat.compiling"
      class="primary--text"
    >编译中</div>
    <div v-else-if="status===$consts.codeStat.accepted"
      class="success--text"
    >已接受</div>
    <div v-else-if="status===$consts.codeStat.complErr"
      class="error--text"
    >编译错误</div>
    <div v-else
      class="error--text"
    >系统错误</div>
  </div>
  <div v-else>
    <div v-if="status===$consts.codeStat.pending"
      class="secondary--text"
    >等待中</div>
    <div v-else-if="status===$consts.codeStat.compiling"
      class="primary--text"
    >处理中</div>
    <div v-else-if="status===$consts.codeStat.accepted"
      class="success--text"
    >已结束</div>
    <div v-else-if="status===$consts.codeStat.complErr"
      class="error--text"
    >编译错误</div>
    <div v-else
      class="error--text"
    >系统错误</div>
  </div>
</template>

<script>
export default {
  props: {
    status: Number,
    mode: String
  }
}
</script>

<style>

</style>

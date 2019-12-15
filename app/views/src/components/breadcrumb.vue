<template>
  <div class="breadcrumb">
    <i class = "el-icon-caret-right" v-if="routeList.length !== 0" style="display: inline-block"></i>
    <el-breadcrumb separator="/" style = "display: inline-block; margin-left: 2px">
      <el-breadcrumb-item v-for="(item) in routeList" :key="item.path">
        <router-link :to="{path: item.path, query: item.query}">{{item.title}}</router-link>
      </el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script>
export default {
  name: 'breadcrumb',
  computed: {
    routeList () {
      return this.$store.state.routeList
    }
  }
}
</script>

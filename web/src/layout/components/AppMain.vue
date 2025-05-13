<template>
    <section class="app-main">
      <router-view v-slot="{ Component }">
        <transition name="fade-transform" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </section>
  </template>
  
  <script setup>
  // 无需额外导入
  </script>
  
  <style scoped>
  .app-main {
    padding: 20px;
    min-height: calc(100vh - 50px);
    position: relative;
    overflow: hidden;
  }
  .fade-transform-enter-active,
  .fade-transform-leave-active {
    transition: all 0.5s;
  }
  .fade-transform-enter-from {
    opacity: 0;
    transform: translateX(-30px);
  }
  .fade-transform-leave-to {
    opacity: 0;
    transform: translateX(30px);
  }
  </style>
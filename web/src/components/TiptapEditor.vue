<template>
  <div v-if="editor">
    <bubble-menu
        class="bubble-menu bg-gray-700 px-3 py-1 rounded-full flex flex-row gap-5 text-white text-xl"
        :tippy-options="{ duration: 100 }"
        :editor="editor"
    >
      <button @click="editor.chain().focus().toggleBold().run()" class="font-bold" :class="{ 'is-active': editor.isActive('bold') }">B</button>
      <button @click="editor.chain().focus().toggleItalic().run()" class="italic" :class="{ 'is-active': editor.isActive('italic') }">I</button>
      <button @click="editor.chain().focus().toggleStrike().run()" class="line-through" :class="{ 'is-active': editor.isActive('strike') }">S</button>
    </bubble-menu>

    <floating-menu
        class="floating-menu bg-gray-700 px-3 py-1 rounded-full flex flex-row gap-5 text-white text-xl"
        :tippy-options="{ duration: 100 }"
        :editor="editor"
    >
      <button @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }">H1</button>
      <button @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }">H2</button>
      <button @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }">H3</button>
      <button @click="editor.chain().focus().toggleHeading({ level: 4 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 4 }) }">H4</button>
      <button @click="editor.chain().focus().toggleBulletList().run()" :class="{ 'is-active': editor.isActive('bulletList') }">List</button>
    </floating-menu>
  </div>

  <editor-content class="editor border-l-2 border-l-slate-300 pl-5 prose m-5 focus:outline-none" :editor="editor" />
</template>

<script>
import {Editor, EditorContent, BubbleMenu, FloatingMenu,} from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

export default {
  name: "TiptapEditor",
  components: {
    EditorContent,
    BubbleMenu,
    FloatingMenu,
  },
  data() {
    return {
      editor: null,
    }
  },
  mounted() {
    this.editor = new Editor({
      content: '<p>🎉</p>',
      extensions: [
        StarterKit,
      ],
    })
  },
  beforeUnmount() {
    this.editor.destroy()
  }
}
</script>

<style scoped>
.bubble-menu button:hover, .bubble-menu button.is-active, .floating-menu button:hover, .floating-menu button.is-active {
  opacity: 1;
}
.editor {
  width: 100% !important;
  max-width: 100% !important;
}
</style>

<template>
  <div class="text-md flex flex-row gap-5" v-if="editor">
    <button @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }">
      H1
    </button>
    <button @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }">
      H2
    </button>
    <button @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }">
      H3
    </button>
    <button class="font-bold" @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }">
     B
    </button>
    <button class="italic" @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }">
      I
    </button>
    <button @click="editor.chain().focus().setTextAlign('left').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'left' }) }">
      <MenuAlt2Icon class="w-4 h-4"></MenuAlt2Icon>
    </button>
    <button @click="editor.chain().focus().setTextAlign('center').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'center' }) }">
      <MenuAlt4Icon class="w-4 h-4"></MenuAlt4Icon>
    </button>
    <button @click="editor.chain().focus().setTextAlign('right').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'right' }) }">
      <MenuAlt3Icon class="w-4 h-4"></MenuAlt3Icon>
    </button>
    <button @click="editor.chain().focus().setTextAlign('justify').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'justify' }) }">
      <MenuIcon class="w-4 h-4"></MenuIcon>
    </button>
    <button @click="setLink" :class="{ 'is-active': editor.isActive('link') }">
      Đặt link
    </button>
    <button @click="editor.chain().focus().unsetLink().run()" :disabled="!editor.isActive('link')">
      Xóa link
    </button>
  </div>

  <editor-content class="editor border-l-2 border-l-slate-300 pl-5 prose m-5 focus:outline-none" :editor="editor" />
</template>

<script>
import {Editor, EditorContent, } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import TextAlign from '@tiptap/extension-text-align'
import {MenuAlt2Icon, MenuAlt3Icon, MenuAlt4Icon, MenuIcon} from "@heroicons/vue/solid";
import Link from '@tiptap/extension-link'

export default {
  name: "TiptapEditor",
  components: {
    EditorContent,
    MenuAlt2Icon,
    MenuAlt3Icon,
    MenuAlt4Icon,
    MenuIcon
  },
  props: {
    content: String
  },
  emits: ['onChange'],
  data() {
    return {
      editor: null,
    }
  },
  mounted() {
    this.editor = new Editor({
      content: this.content,
      extensions: [
        StarterKit,
        TextAlign.configure({
          types: ['heading', 'paragraph'],
        }),
        Link.configure({
          openOnClick: false,
        }),
      ],
      onUpdate: ({ editor }) => {
        const html = editor.getHTML()
        this.$emit('onChange', html)
      },
    })
  },

  methods: {
    setLink() {
      const previousUrl = this.editor.getAttributes('link').href
      const url = window.prompt('URL', previousUrl)
      if (url === null) {
        return
      }
      if (url === '') {
        this.editor
            .chain()
            .focus()
            .extendMarkRange('link')
            .unsetLink()
            .run()

        return
      }
      this.editor
          .chain()
          .focus()
          .extendMarkRange('link')
          .setLink({ href: url })
          .run()
    },
  },
  beforeUnmount() {
    this.editor.destroy()
  }
}
</script>

<style scoped>
.editor {
  width: 100% !important;
  max-width: 100% !important;
}
</style>

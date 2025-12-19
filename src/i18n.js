import { createI18n } from 'vue-i18n';

const messages = {
    en: {
        app: {
            title: 'OneLink',
            subtitle: 'Create your OneLink',
            description: 'All data is stored in the URL. No database needed.',
            createOwn: 'Create your own OneLink',
            livePreview: 'Live Preview',
            note: 'Note: This URL can be long. Use a link shortener if needed!'
        },
        editor: {
            basicInfo: 'Basic Info',
            nameLabel: 'Name',
            namePlaceholder: 'Your Name',
            bioLabel: 'Bio',
            bioPlaceholder: 'Digital Nomad & Creator',
            avatarLabel: 'Avatar URL (Optional)',
            avatarPlaceholder: 'https://...',
            socialPresence: 'Social Presence',
            links: 'Links',
            addLink: 'Add Link',
            linkTitle: 'Link Title',
            generatedLink: 'Your Generated Link',
            copy: 'Copy',
            copied: 'Copied'
        },
        common: {
            anonymous: 'Anonymous',
            previewName: 'Preview Name',
            previewBio: 'This is how your bio will look'
        }
    },
    'zh-CN': {
        app: {
            title: '一链 OneLink',
            subtitle: '创建您的 OneLink',
            description: '所有数据均存储在 URL 中，无需数据库。',
            createOwn: '创建您自己的 OneLink',
            livePreview: '实时预览',
            note: '注意：该 URL 可能较长，建议使用短链接工具！'
        },
        editor: {
            basicInfo: '基本信息',
            nameLabel: '姓名',
            namePlaceholder: '您的姓名',
            bioLabel: '个人简介',
            bioPlaceholder: '数字游民 & 创作者',
            avatarLabel: '头像 URL（可选）',
            avatarPlaceholder: 'https://...',
            socialPresence: '社交媒体',
            links: '自定义链接',
            addLink: '添加链接',
            linkTitle: '链接标题',
            generatedLink: '您生成的链接',
            copy: '复制',
            copied: '已复制'
        },
        common: {
            anonymous: '匿名',
            previewName: '预览姓名',
            previewBio: '这里展示您的个人简介'
        }
    },
    'zh-TW': {
        app: {
            title: '一鏈 OneLink',
            subtitle: '創建您的 OneLink',
            description: '所有數據均存儲在 URL 中，無需數據庫。',
            createOwn: '創建您自己的 OneLink',
            livePreview: '實時預覽',
            note: '注意：該 URL 可能較長，建議使用短鏈接工具！'
        },
        editor: {
            basicInfo: '基本信息',
            nameLabel: '姓名',
            namePlaceholder: '您的姓名',
            bioLabel: '個人簡介',
            bioPlaceholder: '數字遊民 & 創作者',
            avatarLabel: '頭像 URL（可選）',
            avatarPlaceholder: 'https://...',
            socialPresence: '社交媒體',
            links: '自定義鏈接',
            addLink: '添加鏈接',
            linkTitle: '鏈接標題',
            generatedLink: '您生成的鏈接',
            copy: '複製',
            copied: '已複製'
        },
        common: {
            anonymous: '匿名',
            previewName: '預覽姓名',
            previewBio: '這裡展示您的個人簡介'
        }
    },
    ja: {
        app: {
            title: 'ワンリンク OneLink',
            subtitle: 'OneLinkを作成',
            description: 'すべてのデータはURLに保存されます。データベースは不要です。',
            createOwn: '独自のOneLinkを作成する',
            livePreview: 'ライブプレビュー',
            note: '注意：このURLは長くなる可能性があります。リンク短縮サービスの使用をお勧めします！'
        },
        editor: {
            basicInfo: '基本情報',
            nameLabel: '名前',
            namePlaceholder: 'お名前',
            bioLabel: '自己紹介',
            bioPlaceholder: 'デジタルノマド & クリエイター',
            avatarLabel: 'アバターURL（任意）',
            avatarPlaceholder: 'https://...',
            socialPresence: 'ソーシャルメディア',
            links: 'リンク',
            addLink: 'リンクを追加',
            linkTitle: 'タイトル',
            generatedLink: '作成されたリンク',
            copy: 'コピー',
            copied: '完了'
        },
        common: {
            anonymous: '匿名',
            previewName: 'プレビュー名',
            previewBio: 'ここに自己紹介が表示されます'
        }
    }
};

const i18n = createI18n({
    legacy: false,
    locale: navigator.language.startsWith('zh') ? (navigator.language === 'zh-CN' ? 'zh-CN' : 'zh-TW') : (navigator.language.startsWith('ja') ? 'ja' : 'en'),
    fallbackLocale: 'en',
    messages
});

export default i18n;

'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useGoogleLogin, googleLogout } from '@react-oauth/google'

export interface User {
  id: string
  email: string
  name: string
  picture?: string
}

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [idToken, setIdToken] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    // クライアントサイドでのみ実行
    if (typeof window === 'undefined') {
      setLoading(false)
      return
    }

    // ローカルストレージからユーザー情報とトークンを復元
    const storedUser = localStorage.getItem('user')
    const storedToken = localStorage.getItem('id_token')
    
    if (storedUser && storedToken) {
      try {
        setUser(JSON.parse(storedUser))
        setIdToken(storedToken)
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        localStorage.removeItem('user')
        localStorage.removeItem('id_token')
      }
    }
    
    setLoading(false)
  }, [])

  const signInWithGoogle = useGoogleLogin({
    flow: 'auth-code', // Authorization Code Flowを使用してIDトークンを取得
    ux_mode: 'popup', // ポップアップモードを使用（リダイレクトURIの問題を回避）
    onSuccess: async (codeResponse) => {
      try {
        // バックエンドに認証コードを送信してIDトークンを取得
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/v1/auth/google`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ code: codeResponse.code }),
        })

        if (!response.ok) {
          throw new Error('Failed to exchange code for token')
        }

        const data = await response.json()
        
        // IDトークンをデコードしてユーザー情報を取得
        const tokenParts = data.id_token.split('.')
        const payload = JSON.parse(atob(tokenParts[1]))
        
        const user: User = {
          id: payload.sub,
          email: payload.email,
          name: payload.name,
          picture: payload.picture,
        }

        setUser(user)
        setIdToken(data.id_token)
        
        // ローカルストレージに保存
        localStorage.setItem('user', JSON.stringify(user))
        localStorage.setItem('id_token', data.id_token)

        router.push('/input')
      } catch (error) {
        console.error('Error signing in:', error)
      }
    },
    onError: (error) => {
      console.error('Login failed:', error)
    },
  })

  const signOut = () => {
    googleLogout()
    setUser(null)
    setIdToken(null)
    localStorage.removeItem('user')
    localStorage.removeItem('id_token')
    router.push('/')
  }

  return {
    user,
    loading,
    idToken,
    signInWithGoogle,
    signOut,
  }
}

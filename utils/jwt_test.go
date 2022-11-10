package utils

//func TestNewCustomClaims(t *testing.T) {
//	genToken(t)
//}
//
//func genToken(t *testing.T) {
//	testCases := []struct {
//		name         string
//		genTokenFunc func() (string, error)
//		wantRes      any
//		wantErr      error
//	}{
//		{
//			name: "生成token：id=1",
//			genTokenFunc: func() (string, error) {
//				c := NewCustomClaims(1, WithExpire(time.Second))
//				return c.GenToken()
//			},
//			wantRes: "1",
//			wantErr: nil,
//		},
//		{
//			name: "过期",
//			genTokenFunc: func() (string, error) {
//				c := NewCustomClaims(1, WithExpire(time.Second))
//				return c.GenToken()
//			},
//			wantErr: nil,
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := tc.genTokenFunc()
//			assert.Equal(t, err, tc.wantErr)
//			if err != nil {
//				return
//			}
//			assert.Empty(t, res, tc.wantRes)
//		})
//	}
//}
//
//func parseToken(t *testing.T) {
//	testCases := []struct {
//		name         string
//		genTokenFunc func() (string, error)
//		wantRes      any
//		wantErr      error
//	}{
//		{
//			name: "正常解析",
//			genTokenFunc: func() (string, error) {
//				c := NewCustomClaims(1, WithExpire(time.Second))
//				return c.GenToken()
//			},
//			wantRes: "1",
//			wantErr: nil,
//		},
//		{
//			name: "token过期",
//			genTokenFunc: func() (string, error) {
//				c := NewCustomClaims(1, WithExpire(time.Second))
//				token, err := c.GenToken()
//				time.Sleep(time.Second * 2)
//				ParseToken(token)
//			},
//			wantErr: nil,
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := tc.genTokenFunc()
//			assert.Equal(t, err, tc.wantErr)
//			if err != nil {
//				return
//			}
//			assert.Empty(t, res, tc.wantRes)
//		})
//	}
//}

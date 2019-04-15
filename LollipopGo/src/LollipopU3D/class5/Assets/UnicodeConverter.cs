/* ----------------------------------------------------------
文件名称：UnicodeConverter.h

作者：秦建辉

MSN：splashcn@msn.com

当前版本：V1.0

历史版本：
	V1.0	2011年05月12日
			完成正式版本。

功能描述：
	Unicode内码转换器。用于utf-8、utf-16（UCS2）、utf-32（UCS4）之间的编码转换
 ------------------------------------------------------------ */

using System;
using System.Text;

namespace Warfare.Coding
{
    /// <summary>
    /// UTF8、UTF16（UCS2）、UTF32（UCS4）编码转换器
    /// </summary>
    public class UnicodeConverter
    {
        static string base64EncodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
        static int[] base64DecodeChars = new int[]{
            -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
            -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
            -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63,
            52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1,
            -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
            15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1,
            -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
            41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1};

        //base64解码
        public static string base64decode(string str)
        {
            int c1, c2, c3, c4;
            int i, len;//out;
            string out2;
            len = str.Length;
            i = 0;
            out2 = "";
            while (i < len)
            {
                /* c1 */
                do
                {
                    c1 = base64DecodeChars[str[i++] & 0xff];
                } while (i < len && c1 == -1);
                if (c1 == -1)
                    break;
                /* c2 */
                do
                {
                    c2 = base64DecodeChars[str[i++] & 0xff];
                } while (i < len && c2 == -1);
                if (c2 == -1)
                    break;
                out2 += Convert.ToChar((c1 << 2) | ((c2 & 0x30) >> 4));
                /* c3 */
                do
                {
                    c3 = str[i++] & 0xff;
                    if (c3 == 61)
                        return out2;
                    c3 = base64DecodeChars[c3];
                } while (i < len && c3 == -1);
                if (c3 == -1)
                    break;
                out2 += Convert.ToChar(((c2 & 0XF) << 4) | ((c3 & 0x3C) >> 2));
                /* c4 */
                do
                {
                    c4 = str[i++] & 0xff;
                    if (c4 == 61)
                        return out2;
                    c4 = base64DecodeChars[c4];
                } while (i < len && c4 == -1);
                if (c4 == -1)
                    break;
                out2 += Convert.ToChar(((c3 & 0x03) << 6) | c4);
            }
            return out2;
        }

		public static string Base64Encode(string str)
        { //加密 
            string Out = "";
            int i = 0, len = str.Length;
            char c1, c2, c3;
            while (i < len)
            {
                c1 = Convert.ToChar(str[i++] & 0xff);
                if (i == len)
                {
                    Out += base64EncodeChars[c1 >> 2];
                    Out += base64EncodeChars[(c1 & 0x3) << 4];
                    Out += "==";
                    break;
                }
                c2 = str[i++];
                if (i == len)
                {
                    Out += base64EncodeChars[c1 >> 2];
                    Out += base64EncodeChars[((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4)];
                    Out += base64EncodeChars[(c2 & 0xF) << 2];
                    Out += "=";
                    break;
                }
                c3 = str[i++];
                Out += base64EncodeChars[c1 >> 2];
                Out += base64EncodeChars[((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4)];
                Out += base64EncodeChars[((c2 & 0xF) << 2) | ((c3 & 0xC0) >> 6)];
                Out += base64EncodeChars[c3 & 0x3F];
            }
            return Out;
        }


        //utf-16转utf-8
        public static string utf8to16(string str)
        {
            string out2;
            char c;
            int i, len;
            char char2, char3;
            out2 = "";
            len = str.Length;
            i = 0;
            while (i < len)
            {
                c = str[i++];
                switch (c >> 4)
                {
                    case 0:
                    case 1:
                    case 2:
                    case 3:
                    case 4:
                    case 5:
                    case 6:
                    case 7:
                        // 0xxxxxxx
                        out2 += str[i - 1];
                        break;
                    case 12:
                    case 13:
                        // 110x xxxx 10xx xxxx
                        char2 = str[i++];
                        out2 += Convert.ToChar(((c & 0x1F) << 6) | (char2 & 0x3F));
                        break;
                    case 14:
                        // 1110 xxxx 10xx xxxx 10xx xxxx
                        char2 = str[i++];
                        char3 = str[i++];
                        out2 += Convert.ToChar(((c & 0x0F) << 12) |
                            ((char2 & 0x3F) << 6) |
                            ((char3 & 0x3F) << 0));
                        break;
                }
            }
            return out2;
        }

		public string utf16to8(string str) 
		{ 
			string Out = ""; 
			int i, len; 
			char c;//char为16位Unicode字符,范围0~0xffff,感谢vczh提醒 
			len = str.Length; 
			for (i = 0; i < len; i++) 
			{//根据字符的不同范围分别转化 
				c = str[i]; 
				if ((c >= 0x0001) && (c <= 0x007F)) 
				{ 
					Out += str[i]; 
				} 
				else if (c > 0x07FF) 
				{ 
					Out += (char)(0xE0 | ((c >> 12) & 0x0F)); 
					Out += (char)(0x80 | ((c >> 6) & 0x3F)); 
					Out += (char)(0x80 | ((c >> 0) & 0x3F)); 
				} 
				else
				{ 
					Out += (char)(0xC0 | ((c >> 6) & 0x1F)); 
					Out += (char)(0x80 | ((c >> 0) & 0x3F)); 
				} 
			} 
			return Out; 
		} 
    }
}

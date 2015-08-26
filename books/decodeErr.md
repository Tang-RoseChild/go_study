### gob decode: extra data in buffer
### 背景
通过gob，先encode序列化到文件，再通过gob decode，将多个同种类型的map装载到内存中
下面为一部分代码片段：
```
func NewURLStore() *urlStore {
	if store != nil {
		return store
	}
	f, err := os.OpenFile("shorturl.dat", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {

		panic("open file error")
	}

	store := &urlStore{urls: make(map[string]string), file: f}

	if err := store.load(); err != nil {
		fmt.Println("err when load : ", err)
		// panic("load error ")
	}

	return store
}
```
其中load是一个一个的从文件加载到内存中(store).加载完第一个后，再加载第二时，就会出现错误“extra data in buffer”.
* 起初认为是代码问题，但只有一个时，就OK，排除这中case
* 看了源码,如下：

```
func (dec *Decoder) recvMessage() bool {
	// Read a count.
	nbytes, _, err := decodeUintReader(dec.r, dec.countBuf)
	if err != nil {
		dec.err = err
		return false
	}
	if nbytes >= tooBig {
		dec.err = errBadCount
		return false
	}
	dec.readMessage(int(nbytes))
	return dec.err == nil
}

// readMessage reads the next nbytes bytes from the input.
func (dec *Decoder) readMessage(nbytes int) {
	if dec.buf.Len() != 0 {
		// The buffer should always be empty now.
		panic("non-empty decoder buffer")
	}
	// Read the data
	dec.buf.Size(nbytes)
	_, dec.err = io.ReadFull(dec.r, dec.buf.Bytes())
	if dec.err != nil {
		if dec.err == io.EOF {
			dec.err = io.ErrUnexpectedEOF
		}
	}
}

// decodeUintReader reads an encoded unsigned integer from an io.Reader.
// Used only by the Decoder to read the message length.
func decodeUintReader(r io.Reader, buf []byte) (x uint64, width int, err error) {
	width = 1
	n, err := io.ReadFull(r, buf[0:width])
	if n == 0 {
		return
	}
	b := buf[0]
	if b <= 0x7f {
		return uint64(b), width, nil
	}
	n = -int(int8(b))
	if n > uint64Size {
		err = errBadUint
		return
	}
	width, err = io.ReadFull(r, buf[0:n])
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	// Could check that the high byte is zero but it's not worth it.
	for _, b := range buf[0:width] {
		x = x<<8 | uint64(b)
	}
	width++ // +1 for length byte
	return
}

```
>重要的就是上面的三个.可能是调用了read之后，文件的位置发生了变化或者是buffer size发生了变化；具体是怎么回事的话，目前没经历和能力去处理  
* 苦逼的搜索后([连接](http://comments.gmane.org/gmane.comp.lang.go.general/44484))，得到了一个笼统的回答:  
>On Mon, Nov 14, 2011 at 10:08 AM, Rob 'Commander' Pike <r <at> golang.org> wrote:
>The problem is that without special effort you cannot use two
>instantiations of gob.Encoder to write successively to the same file
>and then decode it with gob.Decoder.  
>-rob

#### 最后：
和<< the way to go >>中shorturl类似的一个的练习，将gob改成了json，[链接戳这](https://github.com/nf/goto/commit/16269b8a90b00b7e12331f6632526195b8c3a19f#diff-2)

---------------------------------------------------------------------------------------------------------------------------
崩溃。。。书中"P527--19.7 Using json for storage"有如下描述:
>This is because gob is a stream based protocol that doesn’t support restarting.
解释了上面那个问题。。。

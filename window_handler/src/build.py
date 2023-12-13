import os
from shutil import copy2

from ColorInfo import ColorLogger


class GoBuild:
	def __init__(self, file, name=None, windows=True, linux=True, arm64=True, mips64=True, amd64=True):
		"""

		:param file: 需要构建的主文件(例如: main.go)
		:param name: 生成的执行文件主名称(例如: install)
		:param windows: 是否打包Windows平台
		:param linux: 是否打包Linux平台
		:param arm64: 是否打包arm64架构
		:param mips64: 是否打包mips64架构
		:param amd64: 是否打包amd64架构
		"""
		self.name = name
		self.arch_list = []
		self.os_list = []
		self.amd64 = amd64
		self.mips64 = mips64
		self.arm64 = arm64
		self.linux = linux
		self.windows = windows
		self.file = file
		self.basename = ""
		self.archs = "amd64"
		self.os_type = ""
		self.exe = ""
		self.tmp = ""
		self.logger = ColorLogger()
		self.init()

	def init(self):
		if self.arm64:
			self.arch_list.append("arm64")
		if self.mips64:
			self.arch_list.append("mips64")
		if self.amd64:
			self.arch_list.append("amd64")
		if self.linux:
			self.os_list.append("linux")
		if self.windows:
			self.os_list.append("windows")
		if self.name is None:
			self.basename = str(os.path.basename(self.file)).replace(".go", "")
		else:
			self.basename = self.name

	def delete(self):
		"""
		开始删除生成的临时文件
		:return:
		"""
		tmp = os.path.join(os.getcwd(), self.tmp)
		try:
			os.remove(path=self.tmp)
			self.logger.debug("删除成功: ", tmp)
		except Exception as e:
			self.logger.error(f"删除出错 - [{tmp} ] : ", str(e))

	def copy(self):
		"""
		复制执行文件
		:return:
		"""
		dst = os.path.join("./bin", self.exe)
		self.logger.debug("开始复制: ", dst)
		if os.path.isfile(self.tmp):
			try:
				copy2(src=self.tmp, dst=dst)
				self.delete()
			except Exception as e:
				self.logger.error("复制失败: ", str(e))
		else:
			self.logger.warning("文件不存在: ", self.tmp)

	def build(self):
		self.logger.debug("构建系统: ", self.os_type)
		self.logger.debug("构建架构: ", self.archs)
		self.exe = self.basename + "_" + self.os_type + "-" + self.archs
		self.tmp = str(os.path.basename(self.file)).replace(".go", "")
		if self.os_type == "windows":
			self.exe = self.exe + ".exe"
			self.tmp = self.tmp + ".exe"
		else:
			self.exe = self.exe + ".bin"
		if os.system(f"go build {self.file}") == 0:
			self.logger.info("构建成功,正在生成: ", self.exe)
			self.copy()
		else:
			self.logger.error("构建失败: ", self.exe)

	def ost(self, o):
		os.environ['GOOS'] = o
		self.os_type = o

	def arch(self, arch):
		os.environ['GOARCH'] = arch
		self.archs = arch
		self.build()

	def start(self):
		for i in self.os_list:
			self.ost(i)
			if i == "linux":
				for a in self.arch_list:
					self.arch(arch=a)
			else:
				self.arch(arch="amd64")


if __name__ == "__main__":
	up = GoBuild(file="main.go", name="QNQ")
	up.start()

